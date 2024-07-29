package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/memberlist"
)



func TestHashing(t *testing.T) {
	out := hash("hello world!")
	fmt.Printf("hash value: %v\n", out)
}

func TestSerializeDeserializeKV(t *testing.T) {
	key, value := "test-key", "test-value"
	op := CREATE
	addr := memberlist.Address{
		"hostname:1234",
		"hostname:1234",
	}
	bStr := SerializeKVReq(key, value, op, addr)
	
	kvPair := DeserializeKVReq(bStr)
	if kvPair.Key != key {
		t.Fatalf(`key "%v" not equal to "%v" in deserialized data structure`, key, kvPair.Key)
	}
	if kvPair.Value != value {
		t.Fatalf(`value "%v" not equal to "%v" in deserialized data structure`, value, kvPair.Value)
	}
	if kvPair.Op != op {
		t.Fatalf(`operation "%v" not equal to "%v" in deserialized data structure`, op, kvPair.Op)
	}
	if kvPair.ResNodeAddr != addr {
		t.Fatalf(`Response node addr "%v" not equal to "%v" in deserialized data structure`, addr, kvPair.ResNodeAddr)
	}
}
