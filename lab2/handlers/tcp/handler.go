package tcp

import (
	"bufio"
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

	// incoming request
	scanner := bufio.NewScanner(conn)
	msg := scanner.Text()
	cmd, text, _ := strings.Cut(msg, " ")
	if cmd != "r" && cmd != "w" {
		_, _ = conn.Write([]byte("Invalid command"))
		return
	}

	rndSleep := rand.Intn(7) + 1
	time.Sleep(time.Duration(rndSleep))

	s.mu.Lock()
	defer s.mu.Unlock()

	switch cmd {
	case "r":
		f, err := os.ReadFile(s.FilePath)
		if err != nil {
			log.Print(err)
			_, _ = conn.Write([]byte("Can't read file"))
			return
		}
		_, err = conn.Write(f)
		if err != nil {
			log.Print(err)
			_, _ = conn.Write([]byte("Can't read file"))
			return
		}
	case "w":
		err := os.WriteFile(s.FilePath, []byte(text), 0o666)
		if err != nil {
			log.Print(err)
			_, _ = conn.Write([]byte("Can't write file"))
			return
		}
	}
}
