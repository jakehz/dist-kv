package main

import (
	"fmt"
	"testing"
)

func TestKVStoreGetSetStr(t *testing.T){
	key := "test"
	value := "test"
	kvs := NewKVStore()
	kvs.Set(key, value)
	ret, ok := kvs.Get(key)
	if !ok {
		t.Fatalf(`"%v" key not found in store after setting`, key)
	}
	if ret != value {
		t.Fatalf(`value "%v" at key "%v" not equal to "%v"`, value, key, ret)
	}
}


func TestKVStoreStrNotFound(t *testing.T){
	key := "test"
	kvs := NewKVStore()
	ret, ok := kvs.Get(key)
	if ok {
		t.Fatalf(`"%v" key found in store when not expected`, key)
	}
	if ret != "" {
		t.Fatalf(`value at key "%v" should be EMPTY, not equal to "%v"`, key, ret)
	}
}

func TestGetSerializedUnserialize(t *testing.T) {
	kvs := NewKVStore()
	kvs.Set("test", "test-data")
	kvs.Set("test2", "test2-data")
	b := kvs.GetSerializedMap()
	fmt.Printf("b: %v\n", b)
	kvs.LoadSerializedMap(b)
}
