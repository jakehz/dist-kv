package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
)

type HashRing struct {
	nodes []*memberlist.Node
}

func NewHashRing(capacity int) *HashRing {
	return &HashRing{
		make([]*memberlist.Node, capacity),
	}
}

func (h *HashRing) HashIdx(key string) uint32 {
	return hash(key) % uint32(len(h.nodes))
}

func (h *HashRing) NodeId(node *memberlist.Node) string {
	return fmt.Sprintf("%v:%v", node.FullAddress().Addr, node.Name)
}

func (h *HashRing) GetNode(key string) *memberlist.Node{
	idx := h.HashIdx(key)	
	if len(h.nodes) == 1{
		return h.nodes[0]
	}
	for i := idx; i != idx - 1; i++ {
		if int(i) >= len(h.nodes) {
			i = 0
		}
		if h.nodes[i] != nil {
			return h.nodes[i]
		}
		
	}
	return nil
}

func (h *HashRing) PlaceNode(node *memberlist.Node) {
	if len(h.nodes) == 0{
		return
	}
	id := h.NodeId(node)
	idx := h.HashIdx(id)
	if len(h.nodes) == 1 {
		h.nodes[0] = node
		return
	}
	for i := idx; i != idx - 1; i++ {
		if int(i) >= len(h.nodes){
			i = 0
		}
		if h.nodes[i] == nil {
			h.nodes[i] = node
			return
		}
		
	}
	log.Println("Failed to place node; ring full.")
}

