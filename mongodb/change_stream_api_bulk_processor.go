package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

// 批量操作结构
type bulkOperation struct {
	IndexName string
	DocID     string
	Operation string // index/delete
	Document  map[string]interface{}
	Token     []byte
}

// 批量处理器
type bulkProcessor struct {
	esClient      *elasticsearch.Client
	operations    chan bulkOperation
	maxBatchSize  int
	flushInterval time.Duration
	lastToken     []byte
	tokenFile     string
}

func newBulkProcessor(es *elasticsearch.Client, maxBatchSize int, flushInterval time.Duration, tokenFile string) *bulkProcessor {
	bp := &bulkProcessor{
		esClient:      es,
		operations:    make(chan bulkOperation, 1000),
		maxBatchSize:  maxBatchSize,
		flushInterval: flushInterval,
		tokenFile:     tokenFile,
	}
	go bp.Start()
	return bp
}

func (bp *bulkProcessor) Start() {
	go func() {
		var (
			batch     []bulkOperation
			batchSize int
			ticker    = time.NewTicker(bp.flushInterval)
		)
		defer ticker.Stop()

		for {
			select {
			case op, ok := <-bp.operations:
				if !ok {
					bp.flushBatch(batch)
					return
				}
				batch = append(batch, op)
				batchSize++

				if batchSize >= 1000 {
					bp.flushBatch(batch)
					batch = nil
					batchSize = 0
				}

			case <-ticker.C:
				if batchSize > 0 {
					bp.flushBatch(batch)
					batch = nil
					batchSize = 0
				}
			}
		}
	}()
}

func (bp *bulkProcessor) flushBatch(ops []bulkOperation) {
	if len(ops) == 0 {
		return
	}

	var buf bytes.Buffer
	for _, op := range ops {
		meta := map[string]interface{}{
			"_index": op.IndexName,
			"_id":    op.DocID,
		}

		switch op.Operation {
		case "index":
			buf.WriteString(fmt.Sprintf(`{ "index": %s }%s`, toJSON(meta), "\n"))
			buf.WriteString(toJSON(op.Document) + "\n")
		case "delete":
			buf.WriteString(fmt.Sprintf(`{ "delete": %s }%s`, toJSON(meta), "\n"))
		}
	}
	// 执行批量请求
	res, err := bp.esClient.Bulk(
		&buf,
		bp.esClient.Bulk.WithTimeout(30*time.Second),
		bp.esClient.Bulk.WithRefresh("wait_for"),
	)

	if err == nil && !res.IsError() {
		// 保存最后一个成功的 token
		if len(ops) > 0 {
			saveResumeToken(ops[len(ops)-1].Token)
		}
	} else {
		log.Printf("批量操作失败: %v", err)
		// 这里可以添加重试逻辑
	}
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
