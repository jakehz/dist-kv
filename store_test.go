package main

import (
	"testing"
	"reflect"
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
	kvs1 := NewKVStore()
	kvs2 := NewKVStore()
	key1, val1 := "test", "test-data"
	key2, val2 := "test2", "test2-data"
	kvs1.Set(key1, val1)
	kvs1.Set(key2, val2)
	b := kvs1.GetSerializedMap()
	kvs2.LoadSerializedMap(b)
	v1, ok1 := kvs2.Get(key1)
	if !ok1 {
		t.Fatalf("Key %v not in KV store with loaded serialized data", key1)
	}
	if v1 != val1 {
		t.Fatalf("Value %v does not match value %v in kv store", val1, v1)
	}
	v2, ok2 := kvs2.Get(key2)
	if !ok2 {
		t.Fatalf("Key %v not in KV store with loaded serialized data", key2)
	}
	if v2 != val2 {
		t.Fatalf("Value %v does not match value %v in kv store", val2, v2)
	}
}

func TestSerializeUnserializeMap(t *testing.T){
	data := map[string]string{"test":"val", "test2": "val2"}
	d_bytes := SerializeMap(data)
	d_map := DeserializeMap(d_bytes)
	if !reflect.DeepEqual(data, d_map) {
		t.Fatalf(`map loaded into SerializeMap does not equal map out of Unserialize`)
	}
}
