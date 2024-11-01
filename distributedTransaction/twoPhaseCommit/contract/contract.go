package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var pg *gorm.DB

func init() {
	dsn := "host=localhost user=xianchaoxing password=123456 dbname=contract port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
type ContractDetail struct {
	ID        int    `gorm:"column:id" json:"-"`
	Contract  string `gorm:"column:contract" json:"contract"`
	BeginTime string `gorm:"column:begin_time" json:"begin_time"`
	EndTime   string `gorm:"column:end_time" json:"end_time"`
	// CreatedAt time.Time `gorm:"created" json:"-"`
}

func (b ContractDetail) String() string {
	return "contract_detail"
}

func main() {
	router := gin.Default()
	router.POST("/contracts/dt/:dt", postBilling)
	router.POST("/contracts/commit/dt/:dt", postBillingCommit) // commit or abort
	router.POST("/contracts/abort/dt/:dt", postBillingAbort)   // commit or abort

	router.Run("localhost:9092")
}

// getAlbums responds with the list of all albums as JSON.
func postBilling(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	p := ContractDetail{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	gid := "contract:" + distributedTransactionId
	err = pg.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("contract_detail").Omit("id").Create(&p).Error
		if err != nil {
			return err
		}
		err = tx.Exec("prepare transaction ?", gid).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(500, fmt.Errorf("prepare transaction[%s] failed,%v", gid, err))
		return
	}

	c.Status(200)
}

func postBillingCommit(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	gid := "contract:" + distributedTransactionId
	ret := pg.Exec("commit prepared ?", gid)
	if ret.Error != nil {
		c.AbortWithError(500, fmt.Errorf("commit transaction[%s] failed,%v", gid, ret.Error))
		return
	}

	c.Status(200)
}

func postBillingAbort(c *gin.Context) {
	distributedTransactionId := c.Param("dt")
	gid := "contract:" + distributedTransactionId
	ret := pg.Exec("rollback prepared ?", gid)
	if ret.Error != nil {
		c.AbortWithError(500, fmt.Errorf("commit transaction[%s] failed,%v", gid, ret.Error))
		return
	}

	c.Status(200)
}
