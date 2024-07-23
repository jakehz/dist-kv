package main

import (
	"encoding/gob"
	"bytes"
	"log"

	"github.com/hashicorp/memberlist"
)

type Node struct {
	Name string
	Addr string
	Port int
}

type Cluster struct {
	*memberlist.Memberlist
	LocalNode *Node
	store 	  *KVStore 
}

type KVPair struct {
	Key string `json:"key"`
	Value string `json:"value"`
}


type ClusterDelegate struct{}

func NewCluster(localNode *Node, store *KVStore, logger *log.Logger) (*Cluster, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = localNode.Name
	config.BindAddr = localNode.Addr
	config.BindPort = localNode.Port
	config.Logger = logger
	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	cluster := Cluster{
		Memberlist: list,
		LocalNode:  localNode,
		store:		  store,
	}
	config.Delegate = cluster
	
	return &cluster, nil
}

func (c Cluster) GetNodeForKey(key string) *memberlist.Node{
	// For now we just use a simple implementation
	// replace with Consistent Hashing scheme later down the line
	idx := hash(key) % uint32(c.Memberlist.NumMembers())
	log.Printf(`Sending key "%v" to node %v`, key, idx+1)

	// Memberlist doesn't maintain any specific order; 
	// maybe we can refactor here to hash some attribute of a node 
	// and order it according to the attribute.
	return c.Memberlist.Members()[idx]
}

func (c Cluster) Set(key, value string) {
	// Get the node mapping depending on key
	node := c.GetNodeForKey(key)
	c.Memberlist.SendBestEffort(node, SerializeKV(key, value))
}

// func (c Cluster) SendKVToNode(n *memberlist.Node, key, value string) error {
// 	kv := KVPair{key: key, value: value}
// 	// serialize KV pair
//
// }

func (c Cluster) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)

	return err
}

func (c Cluster) NodeMeta(limit int) []byte{
	return []byte{}
}

func (c Cluster) NotifyMsg(b []byte) {
	log.Println("Recieved KV Pair")
	// Deserialize key value pair
	kv := DeserializeKV(b)
	c.store.Set(kv.Key, kv.Value)
}

func (c Cluster) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{}
}

func (c Cluster) MergeRemoteState(buf []byte, join bool) {
	// Recieves data from remote nodes
	log.Printf("Merge remote state \n")
	// c.store.LoadSerializedMap(buf)
}

func (c Cluster) LocalState(join bool) []byte {
	// Sends local state to remote nodes:
	log.Printf("Sending local data ")
	// return c.store.GetSerializedMap()
	return []byte{}
}

func (c Cluster) pingOtherNode(n *memberlist.Node) {
	msg := []byte("Hello world!")
	c.Memberlist.SendBestEffort(n, msg)
}

func SerializeKV(key, value string) []byte{
	// Similar to SerializeMap... TODO: convert to generics later
	kv := KVPair{key, value}
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(kv)
	if err != nil {
		log.Printf("Error encoding key value pair: %v %v. Error: %v", key, value, err)
		return []byte{}
	}
	return b.Bytes()
}

func DeserializeKV(data []byte) KVPair{
	var kvPair KVPair
	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	
	err := d.Decode(&kvPair)
	if err != nil {
		log.Printf("Error deserializing map: %v", err)
	}
	return kvPair
}
