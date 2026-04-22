package domain

import (
	"reflect"
	"testing"
)

func TestMoney_String(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want string
	}{
		{"£21.70", 2170, "£21.70"},
		{"£21.00", 2100, "£21.0"},
		{"£0.0", 0, "£0.0"},
		{"£0.99", 99, "£0.99"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Money.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMoney(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Money
		wantErr bool
	}{
		{"happey 1", args{"£21.70"}, 2170, false},
		{"happey 2", args{"£21.00"}, 2100, false},
		{"happey 3", args{"£0.0"}, 0, false},
		{"happey 4", args{"£0.99"}, 99, false},
		{"negative major", args{"£-21.70"}, 0, true},
		{"negative minor", args{"£21.-70"}, 0, true},
		{"not pound", args{"21.70$"}, 0, true},
		{"wrong culture", args{"£21,70"}, 0, true},
		{"extra parts", args{"£21.70."}, 0, true},
		{"spaces", args{"£ 21.70"}, 0, true}, // consider allowing spaces where makes sense, like in this example
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoney(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}
