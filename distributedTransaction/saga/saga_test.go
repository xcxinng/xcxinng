package main

import "testing"

func TestSagaService_DoTransactions(t *testing.T) {
	type args struct {
		tp TripParam
	}
	tests := []struct {
		name    string
		s       SagaService
		args    args
		wantErr bool
	}{
		{args: args{
			tp: TripParam{
				Name: "xcx", Destination: "Japan", StartDate: "2024-08-10",
				EndDate: "2024-09-10", PaymentToken: "tkxxkxkxk", Price: "3000RMB"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SagaService{}
			if err := s.DoTransactions(tt.args.tp); (err != nil) != tt.wantErr {
				t.Errorf("SagaService.DoTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
