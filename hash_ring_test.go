package main

import (
	"net"
	"testing"

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

func TestPlaceNodeOne(t *testing.T) {
	//Case: Test where two nodes map to the same index
	const RING_SIZE = 1
	// Create a single node
	hashRing := NewHashRing(RING_SIZE)
	node := NewNode("Node1")
	node1Idx :=  hashRing.HashIdx(hashRing.NodeId(&node))
	hashRing.PlaceNode(&node)
	if hashRing.nodes[node1Idx] == nil && hashRing.nodes[node1Idx].Name == "Node1" {
		t.Fatalf("Failed to place node %v in hash ring.", node)
	}
	// Create another node
	node2 := NewNode("Node2")
	node2Idx :=  hashRing.HashIdx(hashRing.NodeId(&node2))
	if node2Idx != node1Idx{
		t.Fatalf("Node2 idx should map Node1 idx")
	}
	hashRing.PlaceNode(&node2)
	if hashRing.nodes[node2Idx] == nil && hashRing.nodes[node2Idx].Name =="Node2"{
		t.Fatalf("Failed to place node %v in hash ring.", node2)
	}
}

func TestPlaceNodeZero(t *testing.T) {
	// Ensure that placing a node in an empty ring does not result in error
	// or a non-empty ring
	const RING_SIZE = 0
	// Create a single node
	hashRing := NewHashRing(RING_SIZE)
	node := NewNode("Node1")
	// node1Idx :=  hashRing.HashIdx(hashRing.NodeId(&node))
	hashRing.PlaceNode(&node)
	if len(hashRing.nodes) != 0{
		t.Fatalf("Node placed in hash ring of size zero. %v", hashRing.nodes)
	}
}

func TestGetNodeTwo(t *testing.T) {
	/*Case: Test where two nodes map to the same index*/
	const RING_SIZE = 2
	hashRing := NewHashRing(RING_SIZE)
	// Create a single node
	// Node should always map to index 1
	node := NewNode("Node1")
	hashRing.PlaceNode(&node)
	node2 := NewNode("Node2")
	// Create another node that maps to the same idx
	hashRing.PlaceNode(&node2)
	key1 := "key1"
	key2 := "key2"
	resNode1 := hashRing.GetNode(key1)
	if resNode1.FullAddress().Addr != node2.FullAddress().Addr{ 
		t.Fatalf(`Node returned for key %v is not node expected: Expected:"%v" actual:"%v"`, key1, node2.FullAddress().Addr, resNode1.FullAddress().Addr )
	}
	resNode2 := hashRing.GetNode(key2)
	if resNode2.FullAddress().Addr != node.FullAddress().Addr{ 
		t.Fatalf(`Node returned for key %v is not node expected: Expected:"%v" actual:"%v"`, key1, node.FullAddress().Addr, resNode2.FullAddress().Addr )
	}
}

func TestGetNodeOne(t *testing.T) {
	//Case: Test where two nodes map to one node ring
	const RING_SIZE = 1
	hashRing := NewHashRing(RING_SIZE)
	// Create a single node
	// Node should always map to index 1
	node := NewNode("Node1")
	hashRing.PlaceNode(&node)
	node2 := NewNode("Node2")
	// Create another node that maps to the same idx
	hashRing.PlaceNode(&node2)
	key1 := "key1"
	key2 := "key2"
	resNode1 := hashRing.GetNode(key1)
	if resNode1.FullAddress().Addr != node2.FullAddress().Addr{ 
		t.Fatalf(`Node returned for key %v is not node expected: Expected:"%v" actual:"%v"`, key1, node2.FullAddress().Addr, resNode1.FullAddress().Addr )
	}
	resNode2 := hashRing.GetNode(key2)
	if resNode2.FullAddress().Addr != node2.FullAddress().Addr{ 
		t.Fatalf(`Node returned for key %v is not node expected: Expected:"%v" actual:"%v"`, key1, node2.FullAddress().Addr, resNode2.FullAddress().Addr )
	}
}

/*
func TestPlaceNodeZero(t *testing.T) {
	// Ensure that placing a node in an empty ring does not result in error
	// or a non-empty ring
	const RING_SIZE = 0
	// Create a single node
	hashRing := NewHashRing(RING_SIZE)
	node := NewNode("Node1")
	// node1Idx :=  hashRing.HashIdx(hashRing.NodeId(&node))
	hashRing.PlaceNode(&node)
	if len(hashRing.nodes) != 0{
		t.Fatalf("Node placed in hash ring of size zero. %v", hashRing.nodes)
	}
}
*/
