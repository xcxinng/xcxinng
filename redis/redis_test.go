package redis

import "testing"

func TestDeduplication(t *testing.T) {
	runDeduplication()
}

func TestRunHelloWorld(t *testing.T) {
	runHelloWorld()
}

func TestRedLock(t *testing.T) {
	runRedLock()
}
