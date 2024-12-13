package main

import (
	"log"
	"time"
)

type Follower struct {
	Term            int
	Msg             chan Message
	electionTimeout <-chan time.Time
}

func (state *Follower) Run(n *Node) {
	if state.electionTimeout == nil {
		log.Print("Became follower on term ", state.Term)
		state.electionTimeout = time.After(randomElectionTimeout())
	}

	select {
	case <-state.electionTimeout:
		log.Print("Election timeout")
		n.state = &Candidate{Term: state.Term + 1, Msg: state.Msg}
		return
	case msg := <-state.Msg:
		n.state = state.handleMessage(n, msg)
	}
}
func (state *Follower) handleMessage(n *Node, msg Message) state {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Print("Old term, sending update")
		n.sendMessage(msg.Sender(), Message{Type: MessageHeartbeat, Term: state.Term})
		return state
	}
	// A new term, instantly follow.
	if msg.Term > state.Term {
		log.Print("New term, following")
		state.Msg <- msg
		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	switch msg.Type {
	case MessageHeartbeat:
		log.Print("Heartbeat from ", msg.sender)
		state.electionTimeout = time.After(randomElectionTimeout())
		// Another leader was already chosen
	case MessageCandidate:
		log.Print("New candidate, following")
		n.sendMessage(msg.Sender(), Message{
			Type: MessageVote,
			Term: msg.Term,
		})
		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	// Ignore other candidates and other messages
	return state
}
func (state *Follower) ReceiveMessage(msg Message) {
	state.Msg <- msg
}
