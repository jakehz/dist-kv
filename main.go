package main

import (
	"log"
	"os"
	"errors"
	"fmt"
	"strconv"
)

type Config struct {
	ipAddr   string
	nodePort int
	httpPort int
	name     string
}

func main() {
	var config *Config
	var err error
	if len(os.Args) == 1 {
		config, err = parseEnvParams()
	} else {
		config, err = parseParams(os.Args)
	}
	if err != nil {
		log.Fatalf("Error parsing parameters: %v", err)
	}
	log.Printf(`Creating new node "%v" at %v`, config.name, config.FullAddr())

	node := &Node{
		Name: config.name, 
		Addr: config.ipAddr,
		Port: config.nodePort,
	}
	logger := setupLogging(config.name)
	log.SetPrefix(fmt.Sprintf("[%v]", config.name))
	store := NewKVStore()

	// Starts listeners to allow other nodes to join this memberlist.
	cluster, err := NewCluster(node, store, logger)

	if err != nil {
		log.Fatalf("Failed to create cluster %v", err)
	}
	
	api := NewAPI(cluster)

	api.Run(fmt.Sprintf("%v:%v", config.ipAddr, config.httpPort))
}

func setupLogging(nodeName string) *log.Logger{
	prefixLogger := log.New(
		os.Stdout, 
		fmt.Sprintf("[%v]: ", nodeName),
		log.LstdFlags,
		)
	log.SetOutput(prefixLogger.Writer())
	log.SetFlags(0)
	return prefixLogger
}

func parseParams(params []string) (*Config, error){
	if len(params) != 5 {
		return nil, errors.New(
			"Invalid number of params; Enter 3 [name] [IP] [Node Port] [HTTP Port]",
			)
	}
	name     := params[1]
	ipAddr   := params[2]
	nodePort, err := strconv.Atoi(params[3])
	if err != nil {
		log.Fatalf("Error getting node port: %v", err)
	}
	httpPort, err := strconv.Atoi(params[4])
	if err != nil {
		log.Fatalf("Error getting http port: %v", err)
	}

	return &Config{
		name:     name,
		nodePort: nodePort,
		httpPort: httpPort,
		ipAddr:   ipAddr,
	}, nil
}

func parseEnvParams() (*Config, error){
	name     := os.Getenv("NODE_NAME")
	ipAddr   := os.Getenv("IP_ADDR")
	nodePort, err := strconv.Atoi(os.Getenv("NODE_PORT"))
	if err != nil {
		log.Fatalf("Error getting node port: %v", err)
	}
	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatalf("Error getting http port: %v", err)
	}

	return &Config{
		name:     name,
		nodePort: nodePort,
		httpPort: httpPort,
		ipAddr:   ipAddr,
	}, nil
}
func (c *Config) FullAddr() string {
	return fmt.Sprintf("%v:%v", c.ipAddr, c.nodePort)
}
