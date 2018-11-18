package ui

import (
	"Chat/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func Run(m *model.Model, listener net.Listener) {
	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	http.Handle("/scripts/", http.FileServer(http.Dir("")))
	http.Handle("/css/", http.FileServer(http.Dir("")))
	http.Handle("/messages/", messages())

	/// Following 2 actions should be done with GET and POST methods
	/// BTW you can use gorilla/mux and other gorilla's librarie for REST apps. They're really handy
	http.Handle("/messages/get", getMessages(m))
	http.Handle("/messages/send", sendMessage(m))
	go server.Serve(listener)
}

func messages() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "message.html")
	})
}

/// Maybe you wanna get some kind of dependency injection?
func getMessages(m *model.Model) http.Handler { /// Maybe it's better to implement this as method of particular model.Model?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if messages, err := m.Messages(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if stringifyMessages, err := json.Marshal(messages); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if _, err := w.Write(stringifyMessages); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	})
}

func sendMessage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err != nil {
			fmt.Printf("Error! %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			message := model.Message{}
			if err := json.Unmarshal(body, &message); err != nil {
				fmt.Printf("Error! %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if err := m.AddMessage(message); err != nil {
					fmt.Printf("Error! %v", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	})
}
