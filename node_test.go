// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	type args struct {
		perm os.FileMode
		ts   time.Time
		name string
		body string
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{"01", args{0750, now, "", ""}, &Node{0750, now, "", ""}},
		{"02", args{0, now, "==", "==="}, &Node{0, now, "==", "==="}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.perm, tt.args.ts, tt.args.name, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNodeNow(t *testing.T) {
	type args struct {
		name string
		body string
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeNow(tt.args.name, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_SaveAttributes(t *testing.T) {
	type fields struct {
		perm os.FileMode
		time time.Time
		name string
		body string
	}
	type args struct {
		f Fatalfable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				perm: tt.fields.perm,
				time: tt.fields.time,
				name: tt.fields.name,
				body: tt.fields.body,
			}
			n.SaveAttributes(tt.args.f)
		})
	}
}

func TestRfc3339(t *testing.T) {
	type args struct {
		f  Fatalfable
		ts string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rfc3339(tt.args.f, tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rfc3339() = %v, want %v", got, tt.want)
			}
		})
	}
}
