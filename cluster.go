package main
import (
	"github.com/hashicorp/memberlist"
)

type Node struct {
	Name string
	Addr string
}

type Cluster struct {
	*memberlist.Memberlist
	LocalNode *Node
	store 	  *KVStore 
}

func NewCluster(localNode *Node, store *KVStore) (*Cluster, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = localNode.Name
	config.BindAddr = localNode.Addr

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
