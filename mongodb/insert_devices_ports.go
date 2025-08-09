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
	dbName            = "test"
	devicesCollection = "devices"
	portsCollection   = "ports"
	totalDevices      = 20000
	portsPerDevice    = 100  // 200万/2万=100
	batchSize         = 1000 // 批量插入大小
)

var (
	deviceModels = []string{"WS-C3850-48T", "WS-C4500X-32", "Nexus 93180YC-EX", "SRX345"}
	portTypes    = []string{"1G", "10G", "25G", "40G", "100G"}
	statuses     = []int{0, 1, 2}
)

func RunInsert() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// 初始化 MongoDB 客户端
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 获取集合引用
	db := client.Database(dbName)
	devicesColl := db.Collection(devicesCollection)
	portsColl := db.Collection(portsCollection)

	// 清空现有数据
	if err := clearCollections(ctx, devicesColl, portsColl); err != nil {
		log.Fatal(err)
	}

	// 生成设备数据
	deviceIDs, err := generateDevices(ctx, devicesColl)
	if err != nil {
		log.Fatal(err)
	}

	// 生成端口数据
	if err := generatePorts(ctx, portsColl, deviceIDs); err != nil {
		log.Fatal(err)
	}

	fmt.Println("数据插入完成")
}

func clearCollections(ctx context.Context, colls ...*mongo.Collection) error {
	for _, coll := range colls {
		if _, err := coll.DeleteMany(ctx, bson.M{}); err != nil {
			return err
		}
	}
	return nil
}

func generateDevices(ctx context.Context, coll *mongo.Collection) ([]primitive.ObjectID, error) {
	var (
		deviceIDs  = make([]primitive.ObjectID, 0, totalDevices)
		batch      = make([]interface{}, 0, batchSize)
		start      = time.Now()
		lastReport = start
	)

	for i := 1; i <= totalDevices; i++ {
		device := bson.M{
			"name":   fmt.Sprintf("Switch-%06d", i),
			"model":  deviceModels[rand.Intn(len(deviceModels))],
			"ip":     generateIP(),
			"status": statuses[rand.Intn(len(statuses))],
		}
		batch = append(batch, device)

		// 批量插入
		if i%batchSize == 0 || i == totalDevices {
			res, err := coll.InsertMany(ctx, batch)
			if err != nil {
				return nil, err
			}

			// 记录生成的ID
			for _, id := range res.InsertedIDs {
				deviceIDs = append(deviceIDs, id.(primitive.ObjectID))
			}

			batch = batch[:0]
			reportProgress(i, totalDevices, &lastReport, "Devices")
		}
	}

	fmt.Printf("\n生成 %d 台设备，耗时 %v\n", totalDevices, time.Since(start))
	return deviceIDs, nil
}

func generatePorts(ctx context.Context, coll *mongo.Collection, deviceIDs []primitive.ObjectID) error {
	var (
		batch      = make([]interface{}, 0, batchSize)
		counter    = 0
		total      = len(deviceIDs) * portsPerDevice
		start      = time.Now()
		lastReport = start
	)

	for _, deviceID := range deviceIDs {
		for j := 0; j < portsPerDevice; j++ {
			port := bson.M{
				"device_id": deviceID,
				"name":      fmt.Sprintf("Port-%d", j+1),
				"type":      portTypes[rand.Intn(len(portTypes))],
				"speed":     rand.Intn(100000) + 1000, // 1G-100G
				"status":    statuses[rand.Intn(len(statuses))],
				"mac":       generateMAC(),
				"vlan":      rand.Intn(4095) + 1,
			}
			batch = append(batch, port)
			counter++

			if counter%batchSize == 0 || counter == total {
				if _, err := coll.InsertMany(ctx, batch); err != nil {
					return err
				}
				batch = batch[:0]
				reportProgress(counter, total, &lastReport, "Ports")
			}
		}
	}

	fmt.Printf("\n生成 %d 个端口，耗时 %v\n", total, time.Since(start))
	return nil
}

func generateIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(253)+1,
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(253)+1)
}

func generateMAC() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func reportProgress(current, total int, lastReport *time.Time, label string) {
	if time.Since(*lastReport) > 5*time.Second {
		percent := float64(current) / float64(total) * 100
		fmt.Printf("\r%s 进度: %.2f%% (%d/%d)", label, percent, current, total)
		*lastReport = time.Now()
	}
}
