package main

import (
	"encoding/gob"
	"bytes"
	"log"
	"github.com/hashicorp/memberlist"
)

type Operation string
const (
	CREATE Operation = "CREATE"
	READ   Operation = "READ"
	UPDATE Operation = "UPDATE"
	DELETE Operation = "DELETE"
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

type KVReq struct {
	Key         string                `json:"key"`
	Value       string                `json:"value"`
	Op          Operation             `json:"op"`
	ResNodeAddr memberlist.Address    `json:"response_node"`
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
	source_addr := c.Memberlist.LocalNode().FullAddress()
	c.Memberlist.SendBestEffort(
		node, 
		SerializeKVReq(key, value, UPDATE, source_addr),
	)
}

func (c Cluster) Get(key string) {
	/* How do we get a key from another node?
	For now we can just use HTTP. This is probably
	not the most efficient. */
	//idx := c.GetNodeForKey(key)
}
// func (c Cluster) SendKVToNode(n *memberlist.Node, key, value string) error {
// 	kv := KVReq{key: key, value: value}
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
	// This will run on _any_ msg we recieve from another node.
	// We should be able to handle multiple things; key create, read, update, 
	// and delete
	log.Println("Recieved KV Pair")
	// Deserialize key value pair
	kv := DeserializeKVReq(b)
	
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

func SerializeKVReq(key, value string, op Operation, addr memberlist.Address) []byte{
	// Similar to SerializeMap... TODO: convert to generics later
	kv := KVReq{key, value, op, addr}
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(kv)
	if err != nil {
		log.Printf("Error encoding key value pair: %v %v. Error: %v", key, value, err)
		return []byte{}
	}
	return b.Bytes()
}

func DeserializeKVReq(data []byte) KVReq{
	var kvPair KVReq
	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	
	err := d.Decode(&kvPair)
	if err != nil {
		log.Printf("Error deserializing map: %v", err)
	}
	return kvPair
}
