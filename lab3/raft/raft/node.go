package raft

import (
	"fmt"
	"log"
	"net"
	"slices"
	"sync"
)

type Node struct {
	id      int
	state   State
	nodes   []string
	udpConn *net.UDPConn
	mu      sync.Mutex
}

func NewNode(id int, nodes []string, state State) *Node {
	n := &Node{
		id:    id,
		state: state,
		nodes: nodes,
	}

	addr, err := net.ResolveUDPAddr("udp", nodes[id])
	if err != nil {
		log.Fatal("Error resolving UDP address: ", err)
	}

	n.udpConn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Error listening on UDP address: ", err)
	}

	return n
}

func (n *Node) SendMessage(target int, message Message) {
	msg := fmt.Sprintf("%s:%d", message.Type, message.Term)
	// log.Printf("Sending %s to node %d", msg, target)

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

func (n *Node) Broadcast(message Message) {
	for i := range n.nodes {
		if i != n.id {
			n.SendMessage(i, message)
		}
	}
}

func (n *Node) SetState(s State) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = s
}

func (n *Node) ClusterSize() int {
	return len(n.nodes)
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
		// log.Printf("Received %s from %d", msgStr, sender)
		n.mu.Lock()
		n.state.ReceiveMessage(msg)
		n.mu.Unlock()

	}
}

func (n *Node) Run() {
	go n.handleMessages()

	for {
		n.state.Run(n)
	}
}
