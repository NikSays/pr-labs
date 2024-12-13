package state

import (
	"lab3/raft"
	"lab3/timeout"

	"log"
	"time"
)

type Candidate struct {
	Term              int
	Msg               chan raft.Message
	votes             int
	reelectionTimeout <-chan time.Time
}

func (state *Candidate) Run(n *raft.Node) {
	// Initialize state
	if state.votes == 0 {
		log.Print("Became candidate on term ", state.Term)
		state.votes = 1
		state.reelectionTimeout = time.After(timeout.RandomElectionTimeout())
		log.Print("Started election")
		n.Broadcast(raft.Message{
			Type: raft.MessageCandidate,
			Term: state.Term,
		})
	}

	select {
	case <-state.reelectionTimeout:
		log.Print("Reelection timeout")
		n.SetState(&Candidate{Term: state.Term + 1, Msg: state.Msg})
	case msg := <-state.Msg:
		n.SetState(state.handleMessage(n, msg))
	}

}

func (state *Candidate) handleMessage(n *raft.Node, msg raft.Message) raft.State {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Printf("Old term %d from %d, sending update", msg.Term, msg.Sender())
		n.SendMessage(msg.Sender(), raft.Message{Type: raft.MessageUpdate, Term: state.Term})
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
	case raft.MessageHeartbeat, raft.MessageUpdate:
		// Another leader was already chosen
		log.Printf("Heartbeat on same term from %d, following", msg.Sender())
		return &Follower{Term: msg.Term, Msg: state.Msg}
	case raft.MessageVote:
		// Count vote
		state.votes++
		log.Printf("Got vote from %d, now at %d", msg.Sender(), state.votes)
		if state.votes > n.ClusterSize()/2 {
			log.Print("Enough votes for leader")
			return &Leader{Term: state.Term, Msg: state.Msg}
		} else {
			return state
		}
	case raft.MessageCandidate:
		log.Printf("Ignoring another candidate %d on same term", msg.Sender())
	}

	return state
}

func (state *Candidate) ReceiveMessage(msg raft.Message) {
	state.Msg <- msg
}
