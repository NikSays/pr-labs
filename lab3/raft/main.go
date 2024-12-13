package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
	"sync"
)

type state interface {
	Run(ctx *Node)
	ReceiveMessage(msg Message)
}

type Node struct {
	id      int
	state   state
	nodes   []string
	udpConn *net.UDPConn
	mu      sync.Mutex
}

func NewNode(id int, nodes []string) *Node {
	n := &Node{
		id:    id,
		state: &Follower{Term: 0, Msg: make(chan Message, 10)},
		nodes: nodes,
	}

	addr, err := net.ResolveUDPAddr("udp", nodes[id])
	if err != nil {
		panic(err)
	}

	n.udpConn, err = net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	return n
}

func (n *Node) sendMessage(target int, message Message) {
	msg := fmt.Sprintf("%s:%d", message.Type, message.Term)
	// log.Printf("Sending %s to %d", msg, target)

	targetAddr, err := net.ResolveUDPAddr("udp", n.nodes[target])
	if err != nil {
		log.Print("Error resolving target address: ", err)
		return
	}

	_, err = n.udpConn.WriteToUDP([]byte(msg), targetAddr)
	if err != nil {
		log.Print("Error sending message: ", err)
	}
}

func (n *Node) broadcast(message Message) {
	for i := range n.nodes {
		if i != n.id {
			n.sendMessage(i, message)
		}
	}
}

func (n *Node) SetState(s state) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = s
}

func (n *Node) handleMessages() {
	buf := make([]byte, 1024)
	for {
		bytesRead, addr, err := n.udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Print("Error reading UDP message: ", err)
			continue
		}

		msgStr := string(buf[:bytesRead])
		sender := slices.Index(n.nodes, addr.String())
		msg, err := ParseMessage(msgStr, sender)
		if err != nil {
			log.Print("Error parsing message: ", err)
			continue
		}
		// log.Printf("Received %s on term %d from %d", msg.Type, msg.Term, sender)
		n.mu.Lock()
		n.state.ReceiveMessage(msg)
		n.mu.Unlock()

	}
}

func (n *Node) run() {
	go n.handleMessages()

	for {
		n.state.Run(n)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <node_id>")
	}

	nodeID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid node ID")
	}
	log.SetPrefix(fmt.Sprintf("Node %d: ", nodeID))

	nodes := []string{
		"127.0.0.1:8000",
		"127.0.0.1:8001",
		"127.0.0.1:8002",
		"127.0.0.1:8003",
		"127.0.0.1:8004",
		"127.0.0.1:8005",
	}

	node := NewNode(nodeID, nodes)
	defer node.udpConn.Close()

	node.run()
}
