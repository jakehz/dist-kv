package main

import (
	"net"
	"testing"
	//"log"

	"github.com/hashicorp/memberlist"
)

var nodePort uint16 = 12345

func NewNode(name string) memberlist.Node {
	node := memberlist.Node{
		name,
		net.IP{0,0,0,0},
		nodePort,
		[]byte{},
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	}
	nodePort += 1
	return node
}



func TestGetNodeId(t *testing.T){
	hashRing := NewHashRing(2)
	node := NewNode("Node1")
	id := hashRing.NodeId(&node)
	expected := "0.0.0.0:12345:Node1"
	if id != expected{
		t.Fatalf(`Generated invalid id for node. Expected: "%v" Actual: "%v"`, expected,id)
	}
}

func TestHashIdx(t *testing.T) {
	hashRing := NewHashRing(360)
	hashId := "0.0.0.0:12345:Node1"
	idx := hashRing.HashIdx(hashId)
	var expected uint32 = 21
	if idx != expected {
		t.Fatalf("Hashing generated different index than expected: Expected: %v actual: %v", expected, idx)
	}
}

func TestPlaceNodeTwo(t *testing.T) {
	/*Case: Test where two nodes map to the same index*/
	const RING_SIZE = 2
	hashRing := NewHashRing(RING_SIZE)
	// Create a single node
	// Node should always map to index 1
	node := NewNode("Node1")
	node1Idx := hashRing.HashIdx(hashRing.NodeId(&node))
	hashRing.PlaceNode(&node)
	if hashRing.nodes[node1Idx] == nil || hashRing.nodes[node1Idx].Name != "Node1"  {
		t.Fatalf("Failed to place node %v in hash ring.", node)
	}
	node2 := NewNode("Node2")
	// Create another node that maps to the same idx
	node2Idx := hashRing.HashIdx(hashRing.NodeId(&node2))
	if node2Idx != node1Idx {
		t.Fatalf("Node 2 should map to the same index as node1")
	}
	hashRing.PlaceNode(&node2)
	getNode := hashRing.nodes[0]
	if getNode == nil || getNode.Name == "Node2" {
		t.Fatalf("Failed to place node %v in hash ring.", node2)
	}
}
/*
func TestPlaceNodeOne(t *testing.T) {
	//Case: Test where two nodes map to the same index
	const RING_SIZE = 1
	// Create a single node
	val :=  hash("0.0.0.0:12345:Node1") % RING_SIZE 
	node := NewNode("Node1")
	hashRing := NewHashRing(RING_SIZE)
	log.Printf("Place node at index %v", val)
	hashRing.PlaceNode(&node)
	if hashRing.nodes[val] == nil {
		t.Fatalf("Failed to place node %v in hash ring.", node)
	}
	// Create another node
	val2 := hash("0.0.0.0:12346:Node2") % RING_SIZE 
	node2 := NewNode("Node2")
	log.Printf("Place node at index %v", val2)
	hashRing.PlaceNode(&node2)
	if hashRing.nodes[val2] == nil {
		t.Fatalf("Failed to place node %v in hash ring.", node2)
	}
}
*/
