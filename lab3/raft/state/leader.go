package state

import (
	"log"
	"time"

	"lab3/raft"
)

const (
	heartbeatInterval = 2000 * time.Millisecond
)

type Leader struct {
	Term   int
	Msg    chan raft.Message
	ticker *time.Ticker
}

func (state *Leader) Run(n *raft.Node) {
	if state.ticker == nil {
		log.Print("Became leader on term ", state.Term)
		state.ticker = time.NewTicker(heartbeatInterval)
	}

	select {
	case <-state.ticker.C:
		log.Print("Sending heartbeats")
		n.Broadcast(raft.Message{
			Type: raft.MessageHeartbeat,
			Term: state.Term,
		})
	case msg := <-state.Msg:
		n.SetState(state.handleMessage(n, msg))
	}
}

func (state *Leader) handleMessage(n *raft.Node, msg raft.Message) raft.State {
	// Old news, update sender
	if msg.Term < state.Term {
		log.Printf("Old term %d from %d, sending update", msg.Term, msg.Sender())
		n.SendMessage(msg.Sender(), raft.Message{Type: raft.MessageHeartbeat, Term: state.Term})
		return state
	}
	// A new term, instantly follow. Weird to not vote if it's an election,
	// but the simulation does that
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
	case raft.MessageHeartbeat, raft.MessageUpdate, raft.MessageVote:
		log.Printf("Ignoring %s from %d on same term", msg.Type, msg.Sender())

		return &Follower{Term: msg.Term, Msg: state.Msg}
	}

	// Ignore other messages
	return state
}

func (state *Leader) ReceiveMessage(msg raft.Message) {
	state.Msg <- msg
}
