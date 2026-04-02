package domain

import (
	"reflect"
	"testing"
)

func TestNewRateing(t *testing.T) {
	type args struct {
		v uint8
	}
	tests := []struct {
		name    string
		args    args
		want    Rating
		wantErr bool
	}{
		{"in range", args{2}, Rating{2}, false},
		{"above range", args{100}, Rating{}, true},
		{"below range", args{0}, Rating{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRateing(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewRateing() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRating_String(t *testing.T) {
	tests := []struct {
		name string
		r    Rating
		want string
	}{
		{"1", Rating{1}, "*"},
		{"2", Rating{2}, "**"},
		{"3", Rating{3}, "***"},
		{"4", Rating{4}, "****"},
		{"5", Rating{5}, "*****"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("Rating.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
