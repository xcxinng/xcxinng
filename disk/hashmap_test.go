package main

import (
	"testing"
)

func TestFlushToDiskInJSON(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{args: args{m: map[string]interface{}{"hello": "world"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlushToDiskInJSON(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlushToDiskInJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("written bytes: ", got)
		})
	}
}

func TestFlushToDiskInBSON(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{args: args{map[string]interface{}{"hello": "world"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlushToDiskInBSON(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlushToDiskInBSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("written bytes: ", got)
		})
	}
}
