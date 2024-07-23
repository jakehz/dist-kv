package main

import (
	"hash/fnv"
)
func hash(s string) uint32 {
	// create a new Hash object
	h := fnv.New32a()
	// write the string to the hash object as a list of bytes
	h.Write([]byte(s))
	// return the sum (this is the hash value)
	return h.Sum32()
}

