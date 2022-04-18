package rabbitmq

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	send()
}

func TestReceive(t *testing.T) {
	receive()
}

func TestNewTask(t *testing.T) {
	doNewTask(5)
}

func doNewTask(loopCount int) {
	msg := "time-consuming job "
	for i := 0; i < loopCount; i++ {
		rand.Seed(time.Now().UnixNano())
		appendDots := strings.Repeat(".", rand.Intn(5))
		newTask(msg + appendDots)
	}
}

func TestWorker1(t *testing.T) {
	worker()
}

func TestWorker2(t *testing.T) {
	worker()
}
