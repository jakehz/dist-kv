package main

import (
	"log"
)


func main() {
	node := &Node{Name: "node1", Addr: "localhost:1234"}
	store := NewKVStore()
	cluster, err := NewCluster(node, store)

	if err != nil {
		log.Fatalf("Failed to create cluster %v", err)
	}
	
	api := NewAPI(cluster)

	api.Run(":8080")
}
