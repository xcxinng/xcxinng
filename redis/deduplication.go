package redis

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

type OvsFlow struct {
	Id     int    `xorm:"pk autoincr"`
	Flow   string `xorm:"flow"`
	Remark string `xorm:"remark"`
}

func (o OvsFlow) TableName() string {
	return "ovs_flow"
}

// Flow example:
//
// in_port(tap75fc4d97-af),eth(dst=fa:16:3e:8d:c3:24),eth_type(0/0x0800),packets:15008,bytes:90844,used:181.16s,actions:4
//
// in_port(tap75fc4d97-af),eth(src=fa:16:3e:8d:c3:10,dst=fa:16:3e:8d:c3:24),eth_type(0/0x0800),ipv4(src=170.0.0.1,dst=8.8.8.8),udp(dst=67),packets:15008,bytes:90844,used:181.16s,actions:set(tunnel(tun_id=0x2c,dst=7.7.7.7,ttl=64,tp_dst=4789,flags(df|key))),1

type OvsHost interface {
	GetFlows() ([]string, error)
}

type ovsHost struct {
	batch int
}

func newOvsHost(batch int) OvsHost {
	return &ovsHost{batch: batch}
}

func getRandomIP() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(254) + 1
	return fmt.Sprintf("192.168.1.%d", num)
}

// GetFlows implements OvsHost
func (o ovsHost) GetFlows() ([]string, error) {
	flows := make([]string, o.batch)
	for i := 0; i < o.batch; i++ {
		flows[i] = fmt.Sprintf("in_port(tap75fc4d97-af),eth(dst=fa:16:3e:8d:c3:24),eth_type(0/0x0800),ipv4(dst=%s),packets:15008,bytes:90844,used:181.16s,actions:4", getRandomIP())
	}
	return flows, nil
}

type Postgresql interface {
	Insert([]string) (int64, error)
	Close()
}

type myPG struct {
	engine *xorm.Engine
}

func (p myPG) Insert(flows []string) (int64, error) {
	success := 0
	for i, flow := range flows {
		_, err := p.engine.Omit("id").
			Insert(&OvsFlow{Flow: flow, Remark: fmt.Sprintf("remark-%d", i)})
		if err != nil {
			log.Println("[error]:", err, "; flows:", flow)
			continue
		}
		success += 1
	}
	return int64(success), nil
}

func (p myPG) Close() {
	if p.engine != nil {
		p.engine.Close()
	}
}

func newMyPG() (Postgresql, error) {
	engine, err := xorm.NewEngine("postgres",
		"host=localhost user=postgres dbname=postgres  sslmode=disable password=123456 port=5432")
	return &myPG{engine: engine}, err
}

type Redis interface {
	IsFlowExist(string) (bool, error)
	Close()
}

type redisKey struct {
	baseCloser
}

func (r redisKey) IsFlowExist(s string) (bool, error) {
	// key example: 192.168.1.1-{flow_data}
	res, err := r.client.Exists(ctx, r.key(ovsHostIP, s)).Result()
	if err != nil {
		return false, err
	}
	if res == 1 {
		r.client.Expire(ctx, r.key(ovsHostIP, s), time.Minute*2)
		return true, nil
	}
	cmd := r.client.Set(ctx, r.key(ovsHostIP, s), 1, time.Minute*2)
	return false, cmd.Err()
}

type baseCloser struct {
	client *redis.Client
}

func (m baseCloser) Close() {
	if m.client != nil {
		m.client.Close()
	}
}

func (m baseCloser) key(ip, flow string) string {
	return fmt.Sprintf("%s-%s", ip, flow)
}

type redisSet struct {
	baseCloser
}

func (m redisSet) IsFlowExist(s string) (bool, error) {
	exist, err := m.client.SIsMember(ctx, ovsFlowKey, s).Result()
	if err != nil {
		return false, err
	}
	if !exist {
		// insert into ovs_flow_set set
		_, err := m.client.SAdd(ctx, ovsFlowKey, s).Result()
		return false, err
	}
	log.Println("[notice]: flow already exist")
	return true, nil
}

type redisKind uint8

const (
	defaultKind redisKind = iota
	customKind
)

func newMyRedis(k redisKind) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	switch k {
	case customKind:
		return &redisKey{baseCloser{client: rdb}}
	default:
		return &redisSet{baseCloser{client: rdb}}
	}
}

var (
	ctx           = context.Background()
	redisKindFlag = flag.Int("redis_kind", 0, "redis kind")
	keyExpireFlag = flag.Int("key_expire_minute", 2, "redis key expiration time in minute")
)

const (
	ovsFlowKey = "ovs_flow_set"
	ovsHostIP  = "192.168.10.1"
)

func runDeduplication() {
	flag.Parse()
	flows, err := newOvsHost(30).GetFlows()
	if err != nil {
		log.Fatal(err)
	}

	// init redis
	redis := newMyRedis(redisKind(*redisKindFlag))
	defer redis.Close()

	// init pg
	pg, err := newMyPG()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	// process flow data
	var flows2Insert []string
	for _, flow := range flows {
		exist, err := redis.IsFlowExist(flow)
		if err != nil {
			log.Println("[error]:", err, "flow:", flow)
			continue
		}
		if exist {
			log.Println("[existed]:", flow)
			continue
		}
		flows2Insert = append(flows2Insert, flow)
	}
	count, err := pg.Insert(flows2Insert)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert ", count, " flows successfully")
}
