package main

import (
	"log"
	"time"
)

type Candidate struct {
	Term              int
	Msg               chan Message
	votes             int
	reelectionTimeout <-chan time.Time
}

func (state *Candidate) Run(n *Node) {
	// Initialize state
	if state.votes == 0 {
		log.Print("Became candidate on term ", state.Term)
		state.votes = 1
		state.reelectionTimeout = time.After(randomElectionTimeout())
	}

	log.Print("Started election")
	n.broadcast(Message{
		Type: MessageCandidate,
		Term: state.Term,
	})
	select {
	case <-state.reelectionTimeout:
		log.Print("Reelection timeout")
		n.state = &Candidate{Term: state.Term + 1, Msg: state.Msg}
	case msg := <-state.Msg:
		n.state = state.handleMessage(n, msg)
	}

}

func (state *Candidate) handleMessage(n *Node, msg Message) state {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Print("Old term, sending update")
		n.sendMessage(msg.Sender(), Message{Type: MessageUpdate, Term: state.Term})
		return state
	}
	// A new term, instantly follow.
	if msg.Term > state.Term {
		log.Print("New term, following")
		state.Msg <- msg
		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	switch msg.Type {
	case MessageHeartbeat, MessageUpdate:
		// Another leader was already chosen
		log.Print("Heartbeat, following")
		return &Follower{Term: msg.Term, Msg: state.Msg}
	case MessageVote:
		// Count vote
		state.votes++
		log.Print("Got vote, now at ", state.votes)
		if state.votes > len(n.nodes)/2 {
			log.Print("Became leader for term ", state.Term)
			return &Leader{Term: state.Term, Msg: state.Msg}
		} else {
			return state
		}
	case MessageCandidate:
		log.Print("Ignoring another candidate")
	}

	return state
}

func (state *Candidate) ReceiveMessage(msg Message) {
	state.Msg <- msg
}
