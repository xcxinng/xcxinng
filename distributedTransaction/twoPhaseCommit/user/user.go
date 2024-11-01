package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var pg *gorm.DB

func init() {
	dsn := "host=localhost user=xianchaoxing password=123456 dbname=customer port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
type UserDetail struct {
	ID     int    `gorm:"column:id" json:"-"`
	Name   string `gorm:"column:name" json:"name"`
	Gender string `gorm:"column:gender" json:"gender"`
	Age    int    `gorm:"column:age" json:"age"`
}

func (b UserDetail) String() string {
	return "user_detail"
}

/*
transaction | gid   |           prepared            |    owner     | database
-------------+--------------------------------------+--------------------------
        886 | a0a59 | 2024-08-20 18:29:33.783179+08 | xianchaoxing | billing
*/

type preparedTransaction struct {
	Gid string `gorm:"column:gid"`
}

func getTxStatus(id string) (string, error) {
	resp, err := http.Get("http://localhost:9090/services/gid/" + id)
	if err != nil {
		return "nil", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	defer resp.Body.Close()
	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}
	var ret = struct {
		Status string `json:"status"`
	}{}
	err = json.Unmarshal(bytesData, &ret)
	if err != nil {
		return "nil", err
	}
	return ret.Status, nil
}

func genGid(id string) string {
	return "customer:" + id
}

func main() {
	// check unfinished distributed transaction
	go func() {
		var txs []preparedTransaction
		err := pg.Raw("select gid from pg_prepared_xacts where database = 'customer';").Scan(&txs).Error
		if err != nil {
			fmt.Println(err)
		}

		if len(txs) == 0 {
			fmt.Println("no preparing or committing transaction found")
			return
		}

		for _, tx := range txs {
			status, err := getTxStatus(tx.Gid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if status == "committed" {
				err = commitPreparedTransaction(tx.Gid)
				if err != nil {
					fmt.Printf("commit prepared transaction:%s error:%v\n", tx.Gid, err)
				} else {
					fmt.Printf("prepared transaction:%s committed\n", tx.Gid)
				}
			} else {
				err = abortPreparedTransaction(tx.Gid)
				if err != nil {
					fmt.Printf("abort prepared transaction:%s error:%v\n", tx.Gid, err)
				} else {
					fmt.Printf("prepared transaction:%s aborted\n", tx.Gid)
				}
			}
		}
	}()
	router := gin.Default()
	router.POST("/users/dt/:dt", postBilling)
	router.POST("/users/commit/dt/:dt", postBillingCommit) // commit or abort
	router.POST("/users/abort/dt/:dt", postBillingAbort)   // commit or abort

	router.Run("localhost:9091")
}

// getAlbums responds with the list of all albums as JSON.
func postBilling(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	p := UserDetail{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	err = pg.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user_detail").Omit("id").Create(&p).Error
		if err != nil {
			return err
		}
		err = prepareTransaction(distributedTransactionId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("prepare transaction[%s] failed,%v", distributedTransactionId, err))
		return
	}

	c.Status(200)
}

func postBillingCommit(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	err := commitPreparedTransaction(distributedTransactionId)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("commit transaction[%s] failed,%v", distributedTransactionId, err))
		return
	}

	c.Status(200)
}

func commitPreparedTransaction(id string) error {
	gid := genGid(id)
	return pg.Exec("commit prepared ?", gid).Error
}

func abortPreparedTransaction(id string) error {
	gid := genGid(id)
	return pg.Exec("rollback prepared ?", gid).Error
}
func prepareTransaction(id string) error {
	gid := genGid(id)
	return pg.Exec("prepare transaction ?", gid).Error
}

func postBillingAbort(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	err := abortPreparedTransaction(distributedTransactionId)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("commit transaction[%s] failed,%v", distributedTransactionId, err))
		return
	}

	c.Status(200)
}
