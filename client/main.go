package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func mockedIP() string {
	var intArray [4]int
	for i := 0; i < len(intArray); i++ {
		rand.Seed(time.Now().UnixNano())
		intArray[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", intArray[0], intArray[1], intArray[2], intArray[3])
}

func connect(ip string) (ws *websocket.Conn, err error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", ip)
}

// function to send message
func send(username string, text string, ws *websocket.Conn) {
	m := Message{
		Username: username,
		Text:     text,
	}
	err := websocket.JSON.Send(ws, m)
	if err != nil {
		fmt.Println("Could not send message: ", err.Error())
	}
}

var port = flag.String("port", "9000", "Port used for websocket connection.")
var username = flag.String("username", "Unknown", "Username to display.")

func main() {
	flag.Parse()

	ip := mockedIP()
	ws, err := connect(ip)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	send(*username, "Client connected with IP: "+ip, ws)

	// receiving message
	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())
				break
			}
			fmt.Println(m.Username+": ", m.Text)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		} else if text == "exit" {
			ws.Close()
			os.Exit(0)
		} else {
			send(*username, text, ws)
		}

	}
}
