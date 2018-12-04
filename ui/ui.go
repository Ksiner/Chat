package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/Ksiner/Chat/model"

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
	r.HandleFunc("/", authorizeForm)
	r.Handle("/authorize", authorization(m))
	r.HandleFunc("/messages", messages)
	r.Handle("/messages/request", getMessages(m)).Methods("GET")
	r.Handle("/messages/request", sendMessage(m)).Methods("POST")
	r.Handle("/users/request", getUsers(m)).Methods("GET")
	return r
}

func messages(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "message.html")
}

func authorization(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{Name: "username", Value: r.FormValue("username")}
		http.SetCookie(w, cookie)
		ok, err := m.CheckCurrentUser(model.User{Name: cookie.Value})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			if err := m.AddUser(model.User{Name: cookie.Value}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if _, err := w.Write([]byte("messages")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//http.Redirect(w, r, "/messages?username="+r.FormValue("username"), http.StatusFound)
	})
}

func authorizeForm(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("username"); err != nil {
		http.ServeFile(w, r, "authorize.html")
	} else {
		messages(w, r)
	}
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	http.ServeFile(w, r, vars["filegroup"]+"/"+vars["filename"])
}

func checkUnauthorize(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("username")
	if err != nil {
		authorizeForm(w, r)
		return nil
	}
	fmt.Printf("COOKIE!!! %v \n", cookie.Value)
	return cookie
}

func getMessages(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := checkUnauthorize(w, r)
		if cookie != nil {
			messages := make([]*model.Message, 0)
			var err error
			if r.FormValue("targetuser") != "" {
				dialogUsers := make(map[string]interface{}, 0)
				dialogUsers["currentuser"] = cookie.Value
				dialogUsers["targetuser"] = r.FormValue("targetuser")
				if messages, err = m.Messages(dialogUsers); err != nil {
					fmt.Print(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			stringifyMessages := make([]byte, 0)
			if stringifyMessages, err = json.Marshal(messages); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(stringifyMessages); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
}

func getUsers(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := checkUnauthorize(w, r)
		if cookie != nil {
			users, err := m.GetUsers(model.User{Name: cookie.Value})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			stringifyUsers, err := json.Marshal(users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(stringifyUsers); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	})
}

func sendMessage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := checkUnauthorize(w, r)
		if cookie != nil {
			if body, err := ioutil.ReadAll(r.Body); err != nil {
				fmt.Printf("Error! %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				message := model.Message{}
				if err := json.Unmarshal(body, &message); err != nil {
					fmt.Printf("Error! %v", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					message.Userfrom = cookie.Value
					if err := m.AddMessage(message); err != nil {
						fmt.Printf("Error! %v", err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	})
}
