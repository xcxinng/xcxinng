package prometheus

import "testing"

func Test_runInstrumentAppHelloWorld(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "111"},
		{name: "222"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runInstrumentAppHelloWorld()
		})
	}
}
