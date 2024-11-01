package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type SquareNumber struct {
	Number       int    `gorm:"primaryKey,autoIncrement"`
	SquareNumber int    `gorm:"column:squareNumber"`
	Name         string `gorm:"column:name"`
}

// 单条写入和批量写入方式性能对比
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 最低效方式写入
	// start := time.Now()
	// for i := 0; i < 10000; i++ {
	// 	result := db.Table("squareNumber").Create(&SquareNumber{
	// 		Number:       i,
	// 		SquareNumber: i * i,
	// 	}) // Insert tuples (i, i^2)
	// 	if result.Error != nil {
	// 		panic(result.Error) // proper error handling instead of panic in your app
	// 	}
	// }
	// log.Println("timeUseSeconds: ", time.Since(start).Seconds())
	//

	// 批次写入,每个批次在同个事务（应该是driver会做的优化，即使driver没做，我相信InnoDB也会做）
	// 因为批量写入可以减少flush次数，flush是一个同步IO操作，比较耗时
	// 这就是为什么要避免 for...insert 方式写库的原因： 效率/性能问题
	start := time.Now()
	var records []SquareNumber
	for i := 0; i < 10000; i++ {
		s := SquareNumber{
			Number:       i,
			SquareNumber: i * i,
		}
		records = append(records, s)
	}
	result := db.Table("squareNumber").CreateInBatches(&records, 1000) // Insert tuples (i, i^2)
	if result.Error != nil {
		panic(result.Error) // proper error handling instead of panic in your app
	}
	log.Println("timeUseSeconds: ", time.Since(start).Seconds())
	// timeUseSeconds:  0.127933084
}
