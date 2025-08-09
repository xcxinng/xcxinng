package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func PrettyStruct(data interface{}) string {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(val)
}

type JoinConfig struct {
	From         string // 关联集合名
	LocalField   string // 主表关联字段
	ForeignField string // 从表关联字段
	// As           string // 结果字段名
	Filters    bson.M // 子管道过滤条件
	IsOneToOne bool   // 是否1:1关系,否则1:N
}

// BuildJoinPipeline 根据提供的JoinConfig构建pipeline，
// 根collection需要在外层处理，这里不对根表进行任何逻辑处理。
func BuildPipeline(configs []JoinConfig) []bson.M {
	pipeline := make([]bson.M, 0)
	for _, cfg := range configs {
		lookup := bson.M{
			"from":         cfg.From,
			"localField":   cfg.LocalField,
			"foreignField": cfg.ForeignField,
			"as":           cfg.From,
		}

		if len(cfg.Filters) > 0 {
			lookup["pipeline"] = []bson.M{{"$match": cfg.Filters}}
		}

		pipeline = append(pipeline, bson.M{"$lookup": lookup})

		if cfg.IsOneToOne {
			pipeline = append(pipeline, bson.M{
				"$unwind": bson.M{
					"path":                       "$" + cfg.From,
					"preserveNullAndEmptyArrays": true,
				},
			})
		}
	}

	return pipeline
}

func GetMongoClient() *mongo.Client {
	// 初始化MongoDB连接
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type JoinParam struct {
	RootCollection string // 根实体
	RootFilter     bson.M // 根实体过滤条件
	PageSize       int    // page_size和page都是相对根实体而言的
	Page           int
	Config         []JoinConfig // join配置
	Fields         []string     // 需要返回的字段
}

func DynamicJoin(param JoinParam) {
	client := GetMongoClient()
	defer client.Disconnect(context.TODO())

	var pipeline []bson.M
	// 根表处理
	if len(param.RootFilter) > 0 {
		pipeline = append(pipeline, param.RootFilter)
	}
	// 分页尽可能在pipeline前面
	if param.Page > 0 && param.PageSize > 0 {
		pipeline = append(pipeline,
			bson.M{"$skip": (param.Page - 1) * param.PageSize},
			bson.M{"$limit": param.PageSize},
		)
	}
	pipeline = append(pipeline, BuildPipeline(param.Config)...)

	// 添加字段投影
	if len(param.Fields) > 0 {
		pipeline = append(pipeline, buildProjectionStage(param.Fields))
	}
	// fmt.Println(PrettyStruct(pipeline))

	// 执行聚合查询
	ctx := context.TODO()
	op := options.Aggregate()
	//默认aggregate只允许使用100M内存，如果数据集很大，需要用到磁盘空间
	op.SetAllowDiskUse(true)
	//避免无限等待，这种超过分钟级别的肯定有问题的，需要专门优化
	op.SetMaxAwaitTime(time.Minute)
	cursor, err := client.Database("test").
		// 只在secondary节点执行，避免影响写入).
		Collection(param.RootCollection, options.Collection().SetReadPreference(readpref.SecondaryPreferred())).
		Aggregate(ctx, pipeline, op)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	fmt.Println(PrettyStruct(results))
}

func buildProjectionStage(fields []string) bson.M {
	projection := bson.M{"_id": 1} // 默认包含_id

	for _, field := range fields {
		// 处理嵌套字段（如 config.firmware）
		parts := strings.Split(field, ".")
		current := projection

		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = 1
			} else {
				if _, exists := current[part]; !exists {
					current[part] = bson.M{}
				}
				current = current[part].(bson.M)
			}
		}
	}
	return bson.M{"$project": projection}
}

func RunDynamicJoin() {
	param := JoinParam{
		RootCollection: "devices",
		RootFilter:     bson.M{"$match": bson.M{"model": "Nexus 93180YC-EX"}},
		PageSize:       2,
		Page:           1,
		Config: []JoinConfig{
			{
				From:         "ports",
				LocalField:   "_id",
				ForeignField: "device_id",
				Filters: bson.M{"status": 1,
					"name": bson.M{"$regex": primitive.Regex{Pattern: "^Port-6", Options: "i"}}},
			},
			{
				From:         "device_configs",
				LocalField:   "_id",
				ForeignField: "device_id",
				IsOneToOne:   true,
			},
		},
		Fields: []string{"name", "model", "ip",
			"device_configs._id", "device_configs.firmware", "device_configs.last_updated",
			"ports.name", "ports.mac", "ports.vlan",
		},
	}
	DynamicJoin(param)
}
