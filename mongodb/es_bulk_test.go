package main

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

func TestESBulkConnection(t *testing.T) {
	// 1. 创建客户端
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// 如果启用了安全认证需要添加以下配置
		// Username: "elastic",
		// Password: "yourpassword",
		// CACert:   yourCAcert,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	// 2. 检查集群状态
	res, err := es.Info()
	if err != nil {
		t.Fatalf("连接ES失败: %v (请确认ES服务是否运行)", err)
	}
	defer res.Body.Close()
	log.Printf("ES连接成功: %s", res.String())

	// 3. 准备测试数据
	var buf bytes.Buffer
	testDoc := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":    "测试Bulk API连接",
	}

	// 4. 构建Bulk请求体
	meta := map[string]interface{}{
		"index": map[string]interface{}{
			"_index": "connection_test",
			"_id":    "1",
		},
	}
	if err := json.NewEncoder(&buf).Encode(meta); err != nil {
		t.Fatal(err)
	}
	buf.WriteByte('\n')

	if err := json.NewEncoder(&buf).Encode(testDoc); err != nil {
		t.Fatal(err)
	}
	buf.WriteByte('\n')

	// 5. 执行Bulk请求
	res, err = es.Bulk(
		&buf,
		es.Bulk.WithTimeout(5*time.Second),
		es.Bulk.WithPretty(),
		es.Bulk.WithHuman(),
	)
	if err != nil {
		t.Fatalf("Bulk请求失败: %v", err)
	}
	defer res.Body.Close()

	// 6. 解析响应
	if res.IsError() {
		t.Fatalf("ES返回错误: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	// 7. 验证结果
	if result["errors"] == true {
		t.Fatal("批量操作包含错误:", result)
	}

	log.Printf("测试成功！响应: %+v", result)
}
