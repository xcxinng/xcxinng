package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 0. 准备一个replica set集群（单实例没有启用oplog机制不支持change stream api）
// 1. 实现从 MongoDB 监听日志功能
// 2. 实现

func GetClient() *mongo.Client {
	uri := "mongodb://localhost:27017,localhost:27020,localhost:27021/?replicaSet=rs0"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// connect to mongodb
	if err = client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	return client
}

type ChangeEvent struct {
	ID struct {
		Data string `bson:"_data"` // Change Stream 事件唯一标识
	} `bson:"_id"`

	ClusterTime primitive.Timestamp `bson:"clusterTime"` // 操作时间（逻辑时间戳）

	DocumentKey struct {
		ID primitive.ObjectID `bson:"_id"` // 被操作文档的 _id
	} `bson:"documentKey"`

	FullDocument bson.M `bson:"fullDocument"` // 完整文档内容（根据操作类型不同可能不存在）

	Namespace struct {
		Collection string `bson:"coll"` // 集合名称
		Database   string `bson:"db"`   // 数据库名称
	} `bson:"ns"`

	OperationType string    `bson:"operationType"` // 操作类型：insert/update/delete 等
	WallTime      time.Time `bson:"wallTime"`      // 物理服务器时间（带时区）
}

func isResumeTokenExpired(err error) bool {
	return strings.Contains(err.Error(), "resume of change stream was not possible")
}

// 持久化 Resume Token 到文件
func saveResumeToken(token []byte) {
	if err := os.WriteFile("resume_token.bson", token, 0644); err != nil {
		log.Printf("保存恢复令牌失败: %v", err)
	}
}

func startChangeStream(client *mongo.Client) (*mongo.ChangeStream, error) {
	resumeToken := loadResumeToken()

	opts := options.ChangeStream()
	if len(resumeToken) > 0 {
		opts.SetResumeAfter(bson.M{"_data": resumeToken})
		log.Println("从 Resume Token 恢复监听")
	} else {
		log.Println("全新启动，监听所有新事件")
	}

	return client.Database("test").Collection("users").Watch(
		context.Background(),
		mongo.Pipeline{},
		opts,
	)
}

func convertToESOperation(event ChangeEvent) bulkOperation {
	// 生成索引名称 (db_collection)
	indexName := fmt.Sprintf("%s_%s", event.Namespace.Database, event.Namespace.Collection)

	op := bulkOperation{
		IndexName: indexName,
		DocID:     event.DocumentKey.ID.Hex(),
		Token:     []byte(event.ID.Data),
	}

	switch event.OperationType {
	case "insert", "update", "replace":
		op.Operation = "index"
		op.Document = convertDocument(event.FullDocument)
	case "delete":
		op.Operation = "delete"
	}

	return op
}

func convertDocument(doc bson.M) map[string]interface{} {
	// 转换 ObjectID 为字符串
	if id, ok := doc["_id"].(primitive.ObjectID); ok {
		doc["_id"] = id.Hex()
	}

	// 转换日期类型
	for k, v := range doc {
		if dt, ok := v.(primitive.DateTime); ok {
			doc[k] = time.Unix(int64(dt)/1000, 0).UTC().Format(time.RFC3339)
		}
	}
	return doc
}

func createChangeStream(client *mongo.Client, resumeToken []byte) *mongo.ChangeStream {
	opts := options.ChangeStream()
	if len(resumeToken) > 0 {
		opts.SetResumeAfter(bson.M{"_data": string(resumeToken)})
	}

	stream, err := client.Database("test").Watch(
		context.Background(),
		mongo.Pipeline{},
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}
	return stream
}

func connectES() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return es
}

func RunSyncMongoToEs() {
	// 初始化客户端
	mongoClient := GetClient()
	esClient := connectES()
	defer mongoClient.Disconnect(context.Background())

	// 初始化批量处理器
	bp := newBulkProcessor(esClient, 1000, time.Second, "resume_token.bson")
	defer close(bp.operations)

	// 启动同步服务
	sync(mongoClient, bp)
}

const (
	TokenFIle = "resume_token.bson"
)

func loadResumeToken() []byte {
	data, err := os.ReadFile(TokenFIle)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("加载恢复令牌失败: %v", err)
	}
	return data
}

func sync(client *mongo.Client, bp *bulkProcessor) {
	// 加载恢复令牌
	lastToken := loadResumeToken()

	// 创建 Change Stream
	stream := createChangeStream(client, lastToken)
	defer stream.Close(context.Background())

	// 处理事件循环
	for {
		if stream.Next(context.Background()) {
			event := ChangeEvent{}
			err := stream.Decode(&event)
			if err != nil {
				panic(err)
			}

			// 转换为 ES 操作
			op := convertToESOperation(event)
			bp.operations <- op

		} else {
			err := stream.Err()
			if err != nil && isResumeTokenExpired(err) {
				log.Println("Resume Token 已过期，重置监听")
				os.Remove("resume_token.bson")
				stream.Close(context.Background())
				newStream, err := startChangeStream(client)
				if err != nil {
					log.Fatal(err)
				}
				stream = newStream
				continue
			} else {
				log.Fatalf("监听错误: %v", err)
			}
		}
	}

}
