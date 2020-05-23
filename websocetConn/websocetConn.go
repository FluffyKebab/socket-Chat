package websocetConn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"socetChat/cookie"
	"socetChat/storage"
)

type Message struct {
	UserID         string `json:"userID"`
	DisplayName    string `json:"displayName"`
	Body           string `json:"body"`
	IsUsersMessage bool   `json:"isUsersMessage"`
}

type Listener struct {
	ID      string
	Reciver *chan Message
}

type MessageSharer struct {
	Listeners []Listener
	Reciver   chan Message
	Mutex     *sync.Mutex
	Name      string
}

func newMessageSharer(name string) *MessageSharer {
	messageSharer := MessageSharer{
		Listeners: make([]Listener, 0),
		Reciver:   make(chan Message),
		Mutex:     &sync.Mutex{},
		Name:      name,
	}

	return &messageSharer
}

func (m *MessageSharer) addListener(userID string) *Listener {

	reciver := make(chan Message)

	newListener := Listener{
		ID:      userID,
		Reciver: &reciver,
	}

	m.Mutex.Lock()
	m.Listeners = append(m.Listeners, newListener)
	m.Mutex.Unlock()

	return &newListener
}

func (m *MessageSharer) removeListener(listener *Listener) {
	index := -1

	m.Mutex.Lock()

	for i, l := range m.Listeners {
		if l.ID == listener.ID {
			index = i
		}
	}
	if index == -1 {
		fmt.Println("PANIC DETTE SKAL IKKE SKJE, se på remove lister funcen")
	}

	m.Listeners = removeIndex(m.Listeners, index)
	m.Mutex.Unlock()
}

func (m *MessageSharer) listenAndServe() {
	for {
		newMessage := <-m.Reciver

		for _, listener := range m.Listeners {
			go send(listener.Reciver, newMessage)
		}
	}
}

func send(ch *chan Message, messageToSend Message) {
	*ch <- messageToSend
}

var upgrader = websocket.Upgrader{}
var messageSharers map[string]*MessageSharer

func Start() {
	messageSharers = make(map[string]*MessageSharer)
	messageSharers["main"] = newMessageSharer("main")

	for _, sharer := range messageSharers {
		go sharer.listenAndServe()
	}
}

func Handler(storage *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//oppgrader til websocet
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("klarte ikke å oppgradere til en websocet conn", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Finner hvordan chat den skal koble til
		chatNames, ok := r.URL.Query()["name"]
		if !ok || len(chatNames) < 1 {
			fmt.Println("req uten url paramater 'navn'")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chatName := chatNames[0]

		sharer := messageSharers[chatName]
		if sharer == nil {
			if !storage.HasChat(chatName) {
				fmt.Println("Har ikke chat ved det navnet")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			messageSharers[chatName] = newMessageSharer(chatName)
		}

		userSessionID := cookie.GetUserID(r)

		//Legg til i message shareren
		listener := sharer.addListener(userSessionID)

		//HÅNTERER CONN
		done := make(chan bool)
		userData, err := storage.GetUserData(listener.ID)
		if err != nil {
			fmt.Println("Klarte ikke å få brukeren sine data ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//sender alle beskedene fra message shareren til clienten
		go func() {
			for {
				select {
				case message := <-*listener.Reciver:
					message.IsUsersMessage = message.UserID == listener.ID
					message.UserID = ""

					jsonData, err := json.Marshal(message)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Println("Klarte ikke å marshale json fra brruker", err)
						return
					}

					err = conn.WriteMessage(websocket.TextMessage, jsonData)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Println("Klarte ikke å skrive websocet meling", err)
						return
					}

				case <-done:
					return
				}
			}
		}()

		//sender alle beskedene fra clienten til message shareren
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					sharer.removeListener(listener)
					done <- true
					return
				}

				sharer.Reciver <- Message{
					UserID:         listener.ID,
					DisplayName:    userData.DisplayName,
					IsUsersMessage: false,
					Body:           string(message),
				}
			}
		}()
	}
}

func removeIndex(s []Listener, index int) []Listener {
	return append(s[:index], s[index+1:]...)
}
