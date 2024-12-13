package raft

type State interface {
	Run(ctx *Node)
	ReceiveMessage(msg Message)
}
