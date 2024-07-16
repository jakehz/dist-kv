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

	return &Cluster{
		Memberlist: list,
		LocalNode:  localNode,
		store:		  store,
	}, nil
}

func (c *Cluster) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)
	return err
}

func (c *Cluster) NotifyMsg(msg []byte) {

}
