package main

import (
	"log"
	"time"
)

const (
	heartbeatInterval = 2000 * time.Millisecond
)

type Leader struct {
	Term   int
	Msg    chan Message
	ticker *time.Ticker
}

func (state *Leader) Run(n *Node) {
	if state.ticker == nil {
		log.Print("Became leader on term ", state.Term)
		state.ticker = time.NewTicker(heartbeatInterval)
	}

	select {
	case <-state.ticker.C:
		log.Print("Sending heartbeats")
		n.broadcast(Message{
			Type: MessageHeartbeat,
			Term: state.Term,
		})
	case msg := <-state.Msg:
		n.state = state.handleMessage(n, msg)
	}
}

func (state *Leader) handleMessage(n *Node, msg Message) state {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Print("Old term, sending update")
		n.sendMessage(msg.Sender(), Message{Type: MessageHeartbeat, Term: state.Term})
		return state
	}
	// A new term, instantly follow. Weird to not vote if it's an election,
	// but the simulation does that
	if msg.Term > state.Term {
		log.Print("New term, following")
		state.Msg <- msg
		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	switch msg.Type {
	case MessageHeartbeat, MessageUpdate, MessageVote:
		log.Print("Ignoring ", msg.Type)
	case MessageCandidate:
		log.Print("New candidate, following")
		n.sendMessage(msg.Sender(), Message{
			Type: MessageVote,
			Term: msg.Term,
		})
		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	// Ignore other messages
	return state
}

func (state *Leader) ReceiveMessage(msg Message) {
	state.Msg <- msg
}
