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
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error starting server:", err)
		return
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err)
			continue
		}

		go handleNewPlayer(conn)
	}
}

func handleNewPlayer(conn net.Conn) {
	defer conn.Close()
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
			letter, _ := reader.ReadString('\n')
			letter = strings.TrimSpace(letter)
			procesGuess(conn, letter)
		}
	}
}

func setSecretWord(word string) {
	secretWord = word
	guessedWord = make([]byte, len(secretWord))
	for i := range guessedWord {
		guessedWord[i] = '_'
	}
	attempts = len(secretWord) + 5
}

func procesGuess(conn net.Conn, letter string) {
	mu.Lock()
	defer mu.Unlock()
	playerIndex := getPlayerIndex(conn)
	player := &players[playerIndex]

	found := false
	for i, ch := range secretWord {
		if string(ch) == letter {
			guessedWord[i] = byte(ch)
			found = true
			player.score++
		}
	}
	if !found {
		attempts++
	}

	// Broadcast the current guess state
	broadcast(fmt.Sprintf("Word: %s | attempts left: %d\n", string(guessedWord), attempts))
	broadcastScore()

	if strings.Contains(string(guessedWord), "_") && attempts > 0 {
		nextTurn()
	} else if attempts == 0 {
		resetGame()
	} else {
		broadcast("Congrats! the word has been guesses")
		resetGame()
	}

}

// implement getPlayerIndex
func getPlayerIndex(conn net.Conn) int {
	for i, player := range players {
		if player.conn == conn {
			return i
		}
	}
	return -1
}

func nextTurn() {
	playerTurn = (playerTurn + 1) % len(players)
	broadcast(fmt.Sprintf("It's %s turn to guess!\n", players[playerTurn].name))
}

func broadcast(message string) {
	for _, player := range players {
		player.conn.Write([]byte(message))
	}
}

func broadcastScore() {
	message := "Score:\n"
	for _, player := range players {
		message += fmt.Sprintf("%s: %d points\n", player.name, player.score)
	}
}

func resetGame() {
	broadcast("The game is over. Starting new game...\n")
	players = []Player{}
	playerTurn = 0
}
