package main

import (
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	jsonFile = "hashmap.json"
	bsonFile = "hashmap.bson"
)

/*
在应用程序员的角度里，最快能想到以下几种方法将一个hashmap持久化到磁盘中：
（1）JSON序列化(文本文件)
（2）BSON序列化（二进制）
（3）Protobuf（二进制）
但本质上都是一种？
*/
func FlushToDiskInJSON(m map[string]interface{}) (int, error) {
	if len(m) == 0 {
		return 0, nil
	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(jsonFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.WriteAt(jsonBytes, 0)
}

func FlushToDiskInBSON(m map[string]interface{}) (int, error) {
	if len(m) == 0 {
		return 0, nil
	}
	jsonBytes, err := bson.Marshal(m)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(bsonFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.WriteAt(jsonBytes, 0)
}

// func main() {
// 	c, err := FlushToDiskInJSON(map[string]interface{}{"hello": "world"})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("written bytes: ", c)
// }
