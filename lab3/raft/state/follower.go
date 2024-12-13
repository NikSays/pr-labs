package state

import (
	"log"
	"time"

	"lab3/raft"
	"lab3/timeout"
)

type Follower struct {
	Term            int
	Msg             chan raft.Message
	electionTimeout <-chan time.Time
}

func (state *Follower) Run(n *raft.Node) {
	if state.electionTimeout == nil {
		log.Print("Became follower on term ", state.Term)
		state.electionTimeout = time.After(timeout.RandomElectionTimeout())
	}

	select {
	case <-state.electionTimeout:
		log.Print("Election timeout")
		n.SetState(&Candidate{Term: state.Term + 1, Msg: state.Msg})
		return
	case msg := <-state.Msg:
		n.SetState(state.handleMessage(n, msg))
	}
}
func (state *Follower) handleMessage(n *raft.Node, msg raft.Message) raft.State {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Printf("Old term %d from %d, sending update", msg.Term, msg.Sender())
		n.SendMessage(msg.Sender(), raft.Message{Type: raft.MessageHeartbeat, Term: state.Term})
		return state
	}
	// A new term, instantly follow.
	if msg.Term > state.Term {
		log.Printf("New term %d from %d, following", msg.Term, msg.Sender())
		if msg.Type == raft.MessageCandidate {
			log.Print("New candidate, voting")
			n.SendMessage(msg.Sender(), raft.Message{
				Type: raft.MessageVote,
				Term: msg.Term,
			})
		}
		return &Follower{Term: msg.Term, Msg: state.Msg}

	}

	switch msg.Type {
	case raft.MessageHeartbeat:
		log.Print("Heartbeat from ", msg.Sender())
		state.electionTimeout = time.After(timeout.RandomElectionTimeout())
		// Another leader was already chosen
	}

	// Ignore other candidates and other messages
	return state
}
func (state *Follower) ReceiveMessage(msg raft.Message) {
	state.Msg <- msg
}
