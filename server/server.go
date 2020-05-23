package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"socetChat/cookie"
	"socetChat/storage"
	"socetChat/websocetConn"
)

func mainHandler(storage *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Lager en bruker hvis brukeren ikke har det
		if !cookie.HasUser(r) {
			cookieValue, err := storage.NewGuestUser()
			if err != nil {
				fmt.Println(err)
			}

			cookie.Add(w, "sessionID", cookieValue)
		}

		staticServer(w, "./server/src/index.html", "text/html")
	}
}

func Start() {
	websocetConn.Start()

	storage, err := storage.New("./storage/database.db")
	if err != nil {
		fmt.Println("Klarte ikke å lage storage: ", err)
	}

	storage.AddChat("main")

	err = storage.AddChat("main")
	if err != nil {
		fmt.Println("Klarte ikke å lage chat: ", err)
	}

	http.HandleFunc("/", mainHandler(storage))
	http.HandleFunc("/ws/messaageConn", websocetConn.Handler(storage))

	http.HandleFunc("/newChat", newStaticHandler("./server/src/newChat.html", "text/html"))

	http.HandleFunc("/api/usersChats", getUserChatsHandler(storage))
	http.HandleFunc("/api/newestChats", getNewestChatsHandler(storage))
	http.HandleFunc("/api/bigestChats", getBiggestChatsHandler(storage))
	http.HandleFunc("/api/newChat", newChatHandler(storage))

	http.HandleFunc("/css/global", newStaticHandler("./server/src/global.css", "text/css"))
	http.HandleFunc("/css/newChat", newStaticHandler("./server/src/newChat.css", "text/css"))
	http.HandleFunc("/js/chat", newStaticHandler("./server/src/chat.js", "text/js"))
	http.HandleFunc("/js/chooser", newStaticHandler("./server/src/chooser.js", "text/js"))

	http.NotFoundHandler()

	fmt.Println("kjører...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//hjelpe functioner
func sendjsonChats(w http.ResponseWriter, chats []storage.Chat) {
	jsonData, err := json.Marshal(chats)
	if err != nil {
		fmt.Println("Klarte ikke å lage json ut av det ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
}

func newStaticHandler(path, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staticServer(w, path, contentType)
	}
}

func staticServer(w http.ResponseWriter, path, contentType string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Klarte ikke å lese fra fil", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", contentType)
	fmt.Fprint(w, string(data))
}
