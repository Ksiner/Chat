package ui

import (
	"Chat/model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run(m *model.Model, listener net.Listener, cansel context.CancelFunc) {
	server := &http.Server{
		Handler:        makeReouter(m),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go func() {
		defer cansel()
		server.Serve(listener)
	}()
}

func makeReouter(m *model.Model) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{filegroup:(?:css|scripts)}/{filename}", getFiles)
	r.Handle("/messages", messages())
	r.Handle("/messages/request", getMessages(m)).Methods("GET")
	r.Handle("/messages/request", sendMessage(m)).Methods("POST")
	return r
}

func messages() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "message.html")
	})
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	http.ServeFile(w, r, vars["filegroup"]+"/"+vars["filename"])
}

func getMessages(m *model.Model) http.Handler {
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
