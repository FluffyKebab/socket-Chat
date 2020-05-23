package server

import (
	"fmt"
	"net/http"
	"socetChat/cookie"
	"socetChat/storage"
)

func newChatHandler(s *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("feil method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			fmt.Println("Klarte ikke å parse form: ", err)
			w.WriteHeader(500)
			return
		}

		fmt.Println(r.FormValue("name"))

		w.WriteHeader(200)
	}
}

func getUserChatsHandler(s *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := s.GetUserChats(cookie.GetUserID(r))
		if err != nil {
			fmt.Println("klarte ikke å få chats: ", err)
			w.WriteHeader(500)
			return
		}

		sendjsonChats(w, chats)
	}
}

func getNewestChatsHandler(s *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := s.GetNewestChats(0, 20)
		if err != nil {
			fmt.Println("klarte ikke å få chats: ", err)
			w.WriteHeader(500)
			return
		}

		sendjsonChats(w, chats)
	}
}

func getBiggestChatsHandler(s *storage.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := s.GetBiggestChats(0, 20)
		if err != nil {
			fmt.Println("klarte ikke å få chats: ", err)
			w.WriteHeader(500)
			return
		}

		sendjsonChats(w, chats)
	}
}
