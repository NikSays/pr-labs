package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"lab3/raft"
	"lab3/state"
)

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

	node := raft.NewNode(nodeID, nodes, &state.Follower{Term: 0, Msg: make(chan raft.Message, 10)})
	// defer raft.udddpConn.Close()

	node.Run()
}
