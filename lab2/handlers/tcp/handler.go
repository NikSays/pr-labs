package tcp

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	mu       sync.Mutex
	FilePath string
}

func (s *Server) HandleRequest(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	// Skip first newline
	scanner.Scan()
	// Read upto newline
	msg := scanner.Text()
	cmd, text, _ := strings.Cut(msg, " ")
	if cmd != "r" && cmd != "w" {
		_, _ = conn.Write([]byte("Invalid command\n"))
		return
	}

	// Wait for other thread to finish working
	_, err := conn.Write([]byte("Waiting for lock\n"))
	if err != nil {
		log.Print(err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulate delay
	rndSleep := rand.Intn(7) + 1
	_, err = conn.Write([]byte(fmt.Sprintf("Waiting %d seconds\n", rndSleep)))
	if err != nil {
		log.Print(err)
		return
	}
	time.Sleep(time.Duration(rndSleep) * time.Second)

	switch cmd {
	case "r":
		f, err := os.ReadFile(s.FilePath)
		if err != nil {
			log.Print(err)
			_, _ = conn.Write([]byte("Can't read file\n"))
			return
		}
		_, err = conn.Write(append(f, '\n'))
		if err != nil {
			log.Print(err)
			return
		}
	case "w":
		err := os.WriteFile(s.FilePath, []byte(text), 0o666)
		if err != nil {
			log.Print(err)
			_, _ = conn.Write([]byte("Can't write file\n"))
			return
		}
	}

	_, err = conn.Write([]byte("Success\n"))
	if err != nil {
		log.Print(err)
		return
	}
}
