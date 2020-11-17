package main

import (
	"reflect"
	"testing"
)

func Test_checkConnect(t *testing.T) {
	type args struct {
		comment string
	}
	tests := []struct {
		name string
		args args
		want *Struct
	}{
		{"case1", args{comment: "// test test test"}, nil},
		{"case2", args{comment: "//mutagento test test"}, nil},
		{"case3", args{comment: "mutagento testpath test"}, &Struct{Path: "testpath", Name: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkConnect(tt.args.comment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}
