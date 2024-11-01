// Two-Phase Commit(2PC Commit)在不同场景下具体的实现不同，常见领域为数据库和微服务。
//
// 在数据库领域，多见于“server-存储引擎”数据库架构，如没有特别说明，数据库都指mysql.
// 比如server在生成binlog时，如何确保证与存储引擎的WAL一致？格式不要求一致，但两个
// 动作要么发生，要么不发生，必须是原子操作。
//
// 查阅文献
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var pg *gorm.DB

func init() {
	dsn := "host=localhost user=xianchaoxing password=123456 dbname=dt port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)
	pg = db
}

// billing represents data about a record billing.
type SubscriptionDetail struct {
	ID     int    `json:"-"`
	User   string `json:"user"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`

	Bank        string `json:"bank"`
	BankAccount string `json:"bank_account"`
	Amount      int    `json:"amount"`
}

func main() {
	router := gin.Default()
	router.POST("/services/subscription", SubscriptionPOST)
	router.GET("/services/gid/:id", StatusGetHandler)

	router.Run("localhost:9090")
}

func postUser(id string, name string, age int, gender string) error {
	b := map[string]interface{}{
		"age":    age,
		"name":   name,
		"gender": gender,
	}
	dataBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:9091/users/dt/"+id, "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func postBilling(id string, back string, bankAccount string, amount int) error {
	b := map[string]interface{}{
		"account":   bankAccount,
		"from_bank": back,
		"amount":    amount,
	}
	dataBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:9094/bills/dt/"+id, "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func postContract(id string, contract string) error {
	b := map[string]interface{}{
		"contract": contract,
	}
	dataBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:9092/contracts/dt/"+id, "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func StatusGetHandler(c *gin.Context) {
	id := c.Param("id")
	var ret Dt
	err := pg.Table("dt_transaction").Where("id = ?", id).First(&ret).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithError(500, fmt.Errorf("db search failed:%v", err))
		return
	}

	c.JSON(200, map[string]string{"status": ret.Status})
}

// 代表一个类似服务发现之类的组件
var servicesMapping = map[string]string{
	"customer": "http://localhost:9091",
	"contract": "http://localhost:9092",
	"billing":  "http://localhost:9094",
}

// getAlbums responds with the list of all albums as JSON.
func SubscriptionPOST(c *gin.Context) {
	p := SubscriptionDetail{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	id := uuid.New()
	tx := &Dt{Id: id.String(), Status: "preparing"}
	err = pg.Table("dt_transaction").Save(tx).Error
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("saving tx failed: %v", err))
		return
	}

	ch := make(chan int, 3)
	countCh := make(chan int)

	go func() {
		count := 0
		for i := 0; i < 3; i++ {
			value := <-ch
			if value > 0 {
				count++
			}
		}
		countCh <- count
	}()

	go func() {
		err := postUser(id.String(), p.User, p.Age, p.Gender)
		if err == nil {
			ch <- 1
		} else {
			ch <- 0
			fmt.Println(err)
		}
	}()
	go func() {
		err := postBilling(id.String(), p.Bank, p.BankAccount, p.Amount)
		if err == nil {
			ch <- 1
		} else {
			ch <- 0
			fmt.Println(err)
		}
	}()
	go func() {
		err := postContract(id.String(), p.User+":contract")
		if err == nil {
			ch <- 1
		} else {
			ch <- 0
			fmt.Println(err)
		}
	}()
	prepareCount := <-countCh
	if prepareCount != 3 {
		err = abortAll(id.String())
	} else {
		err = commitAll(id.String())
	}
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(500, errors.New(err.Error()))
	}
	c.JSON(200, "success")
}

func abortAll(id string) error {
	resp, err := http.Post("http://localhost:9091/users/abort/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	resp, err = http.Post("http://localhost:9092/contracts/abort/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	resp, err = http.Post("http://localhost:9093/bills/abort/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	err = pg.Table("dt_transaction").Save(&Dt{Id: id, Status: "aborted"}).Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

type Dt struct {
	Id     string `gorm:"column:id"`
	Status string `gorm:"column:status"`
}

func (d Dt) String() string {
	return "dt_transaction"
}

func commitAll(id string) error {
	resp, err := http.Post("http://localhost:9091/users/commit/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	resp, err = http.Post("http://localhost:9092/contracts/commit/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	resp, err = http.Post("http://localhost:9094/bills/commit/dt/"+id, "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	err = pg.Table("dt_transaction").Save(&Dt{Id: id, Status: "committed"}).Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}
