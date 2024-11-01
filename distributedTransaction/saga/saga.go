package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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

type SagaState int8

const (
	SagaStateInit SagaState = iota
	SagaStateBookHotelStart
	SagaStateBookHotelEnd
	SagaStateBookCarStart
	SagaStateBookCarlEnd
	SagaStateStart
	SagaStateEnd
)

type TripParam struct {
	Name         string `json:"name" gorm:"column:name"`
	Destination  string `json:"destination" gorm:"column:destination"`
	StartDate    string `json:"start_date" gorm:"column:start_date"`
	EndDate      string `json:"end_date" gorm:"column:end_date"`
	PaymentToken string `json:"payment_token" gorm:"column:payment_token"`
	Price        string `json:"price" gorm:"column:price"`
}

type CarParam struct {
	SagaId    string
	Model     string
	StartDate string
	EndDate   string
}

type HotelParam struct {
	SagaId      string
	Destination string
	StartDate   string
	EndState    string
}

type TripSaga interface {
	SubTransactions(TripParam) ([]Transaction, error)
}

type Transaction interface {
	Execute(context.Context) error
	Compensation(context.Context) error
	Committed() bool
}

type SagaLog struct {
	Id             string    `gorm:"primaryKey"`
	State          SagaState `gorm:"column:state"`
	Req            string    `gorm:"column:req"`
	HasCompensated bool      `gorm:"column:has_compensated"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (s SagaLog) TableName() string {
	return "saga_log"
}

func save(id string, state SagaState) error {
	return pg.Omit("req").Save(&SagaLog{Id: id, State: state, CreatedAt: time.Now()}).Error
}

func saveCompensated(id string, state SagaState) error {
	return pg.Omit("req").Save(&SagaLog{Id: id,
		State: state, CreatedAt: time.Now(), HasCompensated: true}).Error
}

func insertLog(id string, state SagaState, req TripParam) error {
	jb, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return pg.Create(&SagaLog{Id: id, State: state, Req: string(jb), CreatedAt: time.Now()}).Error
}

func getUnfinishedSaga() ([]SagaLog, error) {
	var ls []SagaLog
	err := pg.Where("created_at > now() - interval '1h' and state != ?", SagaStateEnd).Find(&ls).Error
	return ls, err
}

const (
	jsonType = "application/json"
)

func retryUnfinishedSaga() {
	sgs, err := getUnfinishedSaga()
	if err != nil {
		panic(err)
	}

	saga := SagaService{}
	for _, sg := range sgs {
		if sg.State == SagaStateInit {
			continue
		}

		// requests of Saga must be idempotent
		red := TripParam{}
		err = json.Unmarshal([]byte(sg.Req), &red)
		if err != nil {
			log.Println(err)
			continue
		}

		err = saga.DoTransactions(red)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	// SagaLog用WAL日志实现是最佳方案，这里为了实现简单，用PG代替
	// 进程重启后需要对当前未完成的日志进行继续，必须要等到处理完了才能开始接受新请求
	retryUnfinishedSaga()

	router := gin.Default()
	router.POST("/trips", SubscriptionPOST)
	router.Run("localhost:9090")
}

type SagaService struct{}

func (s SagaService) Redo(id string) error {
	record := TripParam{}
	err := pg.Where("id = ?", id).First(&record).Error
	if err != nil {
		return err
	}
	return s.do(id, record)
}

func (s SagaService) DoTransactions(tp TripParam) error {
	return s.do("", tp)
}

func (s SagaService) do(id string, tp TripParam) error {
	var txs []Transaction
	if id == "" {
		id = uuid.NewString()
	}

	txs = append(txs,
		NewCarTransaction(CarParam{SagaId: id, Model: "宝马三系", StartDate: tp.StartDate, EndDate: tp.EndDate}),
		NewHotelTransaction(HotelParam{SagaId: id, Destination: tp.Destination, StartDate: tp.StartDate, EndState: tp.EndDate}),
	)
	err := insertLog(id, SagaStateStart, tp)
	if err != nil {
		return err
	}

	needCompensation := false
	index := 0
	for i, tx := range txs {
		err := tx.Execute(context.TODO())
		if err != nil {
			needCompensation = true
			index = i
			break
		}
		// if i == 0 {
		// 	return errors.New("saga执行到一半进程崩溃了，此时car已经提交了")
		// }
	}

	if needCompensation {
		fmt.Println("compensating...")
		fmt.Println(">>>>>index: ", index)
		for i := index; i >= 0; i-- {
			if !txs[i].Committed() {
				fmt.Println("index: ", i, "uncommitted, skip compensating")
				continue
			}

			fmt.Println("index: ", i, "committed, need compensating")
			err := txs[i].Compensation(context.TODO())
			if err != nil {
				return err
			}
		}
	}

	if needCompensation {
		err = saveCompensated(id, SagaStateEnd)
	} else {
		err = save(id, SagaStateEnd)
	}
	if err != nil {
		return err
	}

	if needCompensation {
		return errors.New("分布式事务执行失败，操作失败")
	} else {
		return nil
	}
}

type CarService struct {
	currentState SagaState
	param        CarParam
	sagaId       string
}

func (cs *CarService) Execute(ctx context.Context) error {
	cs.currentState = SagaStateBookCarStart
	err := save(cs.sagaId, cs.currentState)
	if err != nil {
		return err
	}
	fmt.Println("book car start")
	err = do("http://localhost:19090/cars/booking", cs.param)
	if err != nil {
		return err
	}

	cs.currentState = SagaStateBookCarlEnd
	err = save(cs.sagaId, cs.currentState)
	fmt.Println("book car end")
	return err
}

func (cs *CarService) Committed() bool {
	return cs.currentState == SagaStateBookCarlEnd
}

func do(url string, i interface{}) error {
	body, err := json.Marshal(i)
	if err != nil {
		return err
	}
	fmt.Println("request url: " + url + ", req body: " + string(body))
	resp, err := http.Post(url, jsonType, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request url:%s error: %v", url, resp.Status)
	}
	return nil
}

func (cs *CarService) Compensation(ctx context.Context) error {
	err := do("http://localhost:19090/cars/booking/undo", cs.param)
	return err
}

func NewCarTransaction(p CarParam) Transaction {
	cs := &CarService{currentState: SagaStateInit, param: p, sagaId: p.SagaId}
	return cs
}

type HotelService struct {
	currentState SagaState
	param        HotelParam
	sagaId       string
}

func (cs *HotelService) Execute(ctx context.Context) error {
	cs.currentState = SagaStateBookHotelStart
	err := save(cs.sagaId, cs.currentState)
	if err != nil {
		return err
	}

	err = do("http://localhost:19091/hotels/booking", cs.param)
	if err != nil {
		return err
	}

	cs.currentState = SagaStateBookHotelEnd
	err = save(cs.sagaId, cs.currentState)
	return err
}

func (cs *HotelService) Committed() bool {
	return cs.currentState == SagaStateBookHotelEnd
}

func (cs *HotelService) Compensation(ctx context.Context) error {
	if cs.currentState != SagaStateBookHotelStart {
		return nil
	}
	err := do("http://localhost:19091/hotels/booking/undo", cs.param)
	return err
}

func NewHotelTransaction(p HotelParam) Transaction {
	cs := &HotelService{currentState: SagaStateInit, param: p, sagaId: p.SagaId}
	return cs
}

func SubscriptionPOST(c *gin.Context) {
	req := TripParam{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	s := SagaService{}
	err = s.DoTransactions(req)
	if err != nil {
		c.AbortWithError(501, err)
	} else {
		c.JSON(200, "success")
	}
}
