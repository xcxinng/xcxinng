package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// dbName            = "network_db"
	// devicesCollection = "devices"
	configsCollection = "device_configs"
	// totalDevices      = 20000 // 复用已有设备数量
	// batchSize         = 1000  // 批量插入大小
)

var (
	firmwareVersions = []string{"15.2(4)E1", "16.9.3", "10.1.2", "8.4.5"}
	snmpVersions     = []string{"v2c", "v3"}
)

func RunInsertDeviceConfig() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// 初始化 MongoDB 客户端
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 获取设备ID列表
	deviceIDs, err := getAllDeviceIDs(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	// 生成设备配置数据
	if err := generateDeviceConfigs(ctx, client, deviceIDs); err != nil {
		log.Fatal(err)
	}

	fmt.Println("1:1 关系数据生成完成")
}

// 获取所有设备ID
func getAllDeviceIDs(ctx context.Context, client *mongo.Client) ([]primitive.ObjectID, error) {
	coll := client.Database(dbName).Collection(devicesCollection)

	cur, err := coll.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var deviceIDs []primitive.ObjectID
	for cur.Next(ctx) {
		var doc struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		deviceIDs = append(deviceIDs, doc.ID)
	}
	return deviceIDs, nil
}

// 生成设备配置数据
func generateDeviceConfigs(ctx context.Context, client *mongo.Client, deviceIDs []primitive.ObjectID) error {
	coll := client.Database(dbName).Collection(configsCollection)

	var (
		batch      = make([]interface{}, 0, batchSize)
		start      = time.Now()
		lastReport = start
	)

	for i, deviceID := range deviceIDs {
		config := bson.M{
			"device_id":    deviceID,
			"firmware":     firmwareVersions[rand.Intn(len(firmwareVersions))],
			"last_updated": time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			"snmp": bson.M{
				"community": fmt.Sprintf("community-%d", rand.Intn(1000)+100),
				"version":   snmpVersions[rand.Intn(len(snmpVersions))],
			},
		}
		batch = append(batch, config)

		// 批量插入
		if (i+1)%batchSize == 0 || i == len(deviceIDs)-1 {
			if _, err := coll.InsertMany(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
			reportInsertConfigProgress(i+1, len(deviceIDs), &lastReport)
		}
	}

	fmt.Printf("\n生成 %d 条设备配置，耗时 %v\n", len(deviceIDs), time.Since(start))
	return nil
}

func reportInsertConfigProgress(current, total int, lastReport *time.Time) {
	if time.Since(*lastReport) > 5*time.Second {
		percent := float64(current) / float64(total) * 100
		fmt.Printf("\r配置数据生成进度: %.2f%% (%d/%d)", percent, current, total)
		*lastReport = time.Now()
	}
}
