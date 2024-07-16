package main

import (
	"log"
	"os"
	"errors"
	"fmt"
)

type Config struct {
	ipAddr string
	port   string
	name   string
}

func main() {
	config, err := parseParams(os.Args)
	if err != nil {
		log.Fatalf("Error parsing parameters: %v", err)
	}
	log.Printf(`Creating new node "%v" at %v`, config.name, config.FullAddr())
	node := &Node{
		Name: config.name, 
		Addr: config.FullAddr(),
	}
	store := NewKVStore()
	cluster, err := NewCluster(node, store)

	if err != nil {
		log.Fatalf("Failed to create cluster %v", err)
	}
	
	api := NewAPI(cluster)

	api.Run(":8080")
}


func parseParams(params []string) (*Config, error){
	if len(params) != 4 {
		return nil, errors.New(
			"Invalid number of params; Enter 3 [name] [IP] [Port]",
			)
	}
	name := params[1]
	ipAddr := params[2]
	port := params[3]
	return &Config{
		name: name,
		port: port,
		ipAddr: ipAddr,
	}, nil
}

func (c *Config) FullAddr() string {
	return fmt.Sprintf("%v:%v", c.ipAddr, c.port)
}
