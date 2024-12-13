package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MessageType string

const (
	MessageHeartbeat MessageType = "heartbeat"
	MessageUpdate    MessageType = "update"
	MessageCandidate MessageType = "candidate"
	MessageVote      MessageType = "vote"
)

type Message struct {
	Type   MessageType
	Term   int
	sender int
}

func (m Message) Sender() int {
	return m.sender
}

func ParseMessage(str string, sender int) (Message, error) {
	typeStr, termStr, ok := strings.Cut(str, ":")

	if !ok {
		return Message{}, errors.New("no separator")
	}

	term, err := strconv.Atoi(termStr)
	if err != nil {
		return Message{}, fmt.Errorf("term is not an int: %w", err)
	}

	return Message{
		Type:   MessageType(typeStr),
		Term:   term,
		sender: sender,
	}, nil
}
