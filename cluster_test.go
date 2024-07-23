package main

import (
	"testing"
	"fmt"
)



func TestHashing(t *testing.T) {
	out := hash("hello world!")
	fmt.Printf("hash value: %v\n", out)
}

func TestSerializeDeserializeKV(t *testing.T) {
	key, value := "test-key", "test-value"
	bStr := SerializeKV(key, value)
	
	kvPair := DeserializeKV(bStr)
	if kvPair.Key != key {
		t.Fatalf(`key "%v" not equal to "%v" in deserialized data structure`, key, kvPair.Key)
	}
	if kvPair.Key != key {
		t.Fatalf(`value "%v" not equal to "%v" in deserialized data structure`, value, kvPair.Value)
	}
}
