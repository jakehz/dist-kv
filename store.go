package main

import (
	"sync"
	"encoding/gob"
	"bytes"
	"log"
)

type KVStore struct {
	data sync.Map
}

func NewKVStore() *KVStore {
	return &KVStore {
		data: sync.Map{},
	}
}

func (s *KVStore) Set(key string, value string) {
	s.data.Store(key, value)
}

func (s *KVStore) Get(key string) (string, bool) {
	value, ok := s.data.Load(key)
	if !ok {
		return "", false
	}

	return value.(string), true
}

func (s *KVStore) Delete(key string) {
	s.data.Delete(key)
}

func (s *KVStore) GetSerializedMap() []byte{
	// TODO: Refactor. Very costly if KV Store is large
	data := make(map[string]string)
	iter := func(key, value any) bool {
		data[key.(string)] = value.(string)
		return true
	}
	s.data.Range(iter)
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	// Encoding the map
 	err := e.Encode(data)
	if err != nil {
			panic(err)
	}
	return b.Bytes()
}

func (s *KVStore) LoadSerializedMap(data []byte) {
	b := bytes.NewBuffer(data)
	var decodedMap map[string]string
	d := gob.NewDecoder(b)

	// Decoding the serialized data
	err := d.Decode(&decodedMap)
	if err != nil {
			panic(err)
	}

	log.Printf("%v", decodedMap)
}


