package intelligentnetwork

import (
	"github.com/IBM/sarama"
)

type NgCmdbREST interface {
	GetInstances(entity, deviceId string) ([]interface{}, error)
}

type Encoder interface {
	NgCmdbREST
	Deserialize(sarama.Message) interface{}
	Resource() string
	Compare([]interface{}) (add, delete, err error)
}

type Arp struct {
}

func (a Arp) GetInstances(entity, deviceId string) ([]interface{}, error) {
	body := map[string]interface{}{}
	_ = body
	// http.Post("switch", "application/json")
	return nil, nil
}

type Lldp struct {
}

type RouteTable struct {
}

type Port struct {
}
