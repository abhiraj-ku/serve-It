package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Player struct {
	conn   net.Conn
	name   string
	score  int
	active bool
}

var players []Player
var playerTurn int
var secretWord string
var guessedWord []byte
var attempts int

var mu sync.Mutex

func main() {
	fmt.Println("Server listening on port :8080")
	server, _ := net.Listen("tcp", ":8080")

	for {
		conn, _ := server.Accept()

		go handleNewPlayer(conn)
	}
}

func handleNewPlayer(conn net.Conn) {
	reader := bufio.NewReader(conn)

	//Ask for player's name
	conn.Write([]byte("enter your name"))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// add this player to players list
	player := Player{
		conn:   conn,
		name:   name,
		score:  0,
		active: false,
	}
	mu.Lock()
	players = append(players, player)
	mu.Unlock()

	if len(players) == 1 {
		player.active = true
		conn.Write([]byte("You are the host! set the secret word: "))
		word, _ := reader.ReadString('\n')
		word = strings.TrimSpace(word)
		mu.Lock()
		setSecretWord(word)
		mu.Unlock()
		broadcast("the word has been set! let the game begin")
		nextTurn()
	} else {
		conn.Write([]byte("waiting for your turn to guess...\n"))
	}

	// start listening to player's command
	for {
		if playerTurn == getPlayerIndex(conn) {
			conn.Write([]byte("Your turn! Guess a letter:"))
			letter, _ := reader.ReadSlice('\n')
			letter = strings.TrimSpace(letter)
			procesguess(conn, letter)
		}
	}
}
