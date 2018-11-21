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
	http.Handle("/messages/users")
	http.Handle("/messages/get", getMessages(m))
	http.Handle("/messages/send", sendMessage(m))
	go server.Serve(listener)
}

func messages() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "message.html")
	})
}

func getUsers(m *model.Model) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if body,err:= ioutil.ReadAll(r.Body);err!=nil{
			fmt.Printf("Error! %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else{
			currentUser:= *model.User{}
			if err:= json.Unmarshal(body,&currentUser);err!=nil{
				fmt.Printf("Error! %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else{
				if users,err:= m.GetUsers(currentUser);err!=nil{
					fmt.Printf("Error! %v", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else{
					if usersJSON,err:=json.Marshal(users);err!=nil{
						fmt.Printf("Error! %v", err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else{
						if _,err:=w.Write(usersJSON);err!=nil{
							fmt.Printf("Error! %v", err.Error())
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
					}
				}
			}
		}
		
	})
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
