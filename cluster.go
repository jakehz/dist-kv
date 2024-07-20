package main
import (
	"github.com/hashicorp/memberlist"
	"log"
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

func (c Cluster) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)
	return err
}

func (c Cluster) NodeMeta(limit int) []byte{
	return []byte{}
}

func (c Cluster) NotifyMsg(b []byte) {
	log.Println("Notify Message")
	log.Printf("%v\n", b)
}

func (c Cluster) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{}
}

func (c Cluster) MergeRemoteState(buf []byte, join bool) {
	// Recieves data from remote nodes
	log.Printf("Merge remote state \n")
	log.Printf("%v \n", buf)
}

func (c Cluster) LocalState(join bool) []byte {
	// Sends local state to remote nodes:
	log.Printf("Sending local data ")
	val, ok := c.store.Get("test")
	if ok {
		return []byte(val)
	} else {
		return []byte{}
	}
}


