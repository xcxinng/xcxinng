// whatever.go only stores those code that does not fit in any other categories
// in this project currently.
//
// PLEASE put your code into the corresponding category as much as possible.

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func testSpilt() {
	// Get the current month as a number
	currentMonth := int(time.Now().Year())

	// Print the current month
	fmt.Println("Current Month (numeric):", currentMonth)
}

type MyStruct struct {
	FloatField float64 `validate:"min=1.0,max=10.1" message:"超出最大最小值范围"`
	IPField    string  `validate:"ip" message:"必须为IP地址格式"`
	CidrField  string  `validate:"cidr"`
	Date       string  `validate:"datetime=2006-01-02"`
	Datetime   string  `validate:"datetime=2006-01-02 15:04:05"`
}

func jsonString() {
	tmp := MyStruct{IPField: "1.1.1.1"}
	fmt.Println(json.Marshal(tmp))
}

func main() {
	// testSpilt()
	// jsonString()

	runHttpPrinter()

	// Server()
	// check()

	// v := validator.New()
	// data := map[string]interface{}{
	// 	"name":     "John",
	// 	"age":      20,
	// 	"ip":       "1.1.1.1",
	// 	"date":     "2020-10-01",
	// 	"datetime": "2020-10-01 15:00:00",
	// }
	// myRules := map[string]interface{}{
	// 	"name":     "lt=10",
	// 	"age":      "lt=100",
	// 	"ip":       "ip",
	// 	"date":     "datetime=2006-01-02",
	// 	"datetime": "datetime=2006-01-02 15:04:05",
	// 	"xcx":      "",
	// }

	// err := v.ValidateMap(data, myRules)
	// for _, e := range err {
	// 	myErr := e.(validator.ValidationErrors)
	// 	fmt.Println(myErr)
	// }
}

func check() {

	// Your date in "2006-12-01" format
	dateStr := "2006-12-01"

	// Parse the date string into a Go time.Time object
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Fatal(err)
	}

	// Format the time.Time as an ISODate-like string
	isoDateString := date.Format("2006-01-02T15:04:05.000Z")
	fmt.Println(isoDateString)
}

type AttributeParam struct {
	EntityID      string `json:"entity_id" binding:"required" uri:"entity_id"`
	AttributeName string `json:"attribute_name" binding:"required"`
}

func Server() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.POST("/entities/:entity_id/attributes", func(c *gin.Context) {
		p := AttributeParam{}
		err := c.ShouldBind(&p)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		err = c.ShouldBindJSON(&p)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		data, err := json.Marshal(p)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		fmt.Println(string(data))
		c.String(http.StatusOK, string(data))
	})
	r.Run("")
	listener, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: "/tmp/gin.sock",
		Net:  "unix",
	})
	if err != nil {
		log.Fatal("Error listening on Unix domain socket:", err)
	}
	if err := http.Serve(listener, r); err != nil {
		log.Fatal("Error serving on Unix domain socket:", err)
	}
}

type Struct1 struct {
	Field1 int
	Field2 string
}

type Struct2 struct {
	Field1 float64
	Field2 string
}

func calculateMD5(data interface{}) string {
	h := md5.New()
	// 使用 fmt.Sprintf 将字段值转换为字符串，并写入哈希
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return hex.EncodeToString(h.Sum(nil))
}

func runMd5() {
	struct1 := Struct1{Field1: 1, Field2: "example"}
	struct2 := Struct2{Field1: 1, Field2: "example"}

	md5sum1 := calculateMD5(struct1)
	md5sum2 := calculateMD5(struct2)

	fmt.Println("MD5 checksum for struct1:", md5sum1)
	fmt.Println("MD5 checksum for struct2:", md5sum2)

	if md5sum1 == md5sum2 {
		fmt.Println("MD5 checksums are equal")
	} else {
		fmt.Println("MD5 checksums are not equal")
	}
}

// 分享题目
// 1. on-disk data structure
//    - B-Tree
//    - LSM_Tree
//
// 2. memory sharing
//    - virtual memory address space
//    - what's fd and what are behind them
//      (1) task_struct + file system + open_files concepts
//      (2) what does a fd mean to a process, not just process-limited, can be shared as well
//      (3) how to get a fd?
//          - open(): used by the initiator
//          - memfd(): used by the initiator
//          - unix domain socket: used by the receiver process
//    - how to share?
//      (1) file-based mem sharing: open() + mmap() + munmap() // flock() if necessary
//      (2) anonymous mem sharing: memfd() + ftruncate() + mmap() + munmap()
//      (3) shm_open()/shm_get() + ftruncate() + shm_at() + shm_ulink()
//    - how to achieve these in Go?
//      (1) vpp stat segment
//      (2) elb listeners health check stat segment
