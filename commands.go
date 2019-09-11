package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type responseErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type User struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	UserName string `json:"username"`
	// Do we want an md5-ed version here?
	Password string `json:"password"`
	Email string `json:"email"`
	Country string `json:"country"`
}

// We mock the database here because it's quite irrelevant
// in this application.
var usersIndex map[string]*User

func InitUsers() {
	usersIndex = make(map[string]*User)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Users_microservice 0.1\n")
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		replyError(http.StatusBadRequest, w, r, err.Error())
		return
	}
	defer r.Body.Close()

	if (userExists(user.UserName)) {
		replyError(http.StatusBadRequest, w, r, "Bad Request")
		return
	}

	usersIndex[user.UserName] = &user

	replyOk(w)
}

func UserEdit(w http.ResponseWriter, r *http.Request) {
	userName := mux.Vars(r)["user"]

	if (!userExists(userName)) {
		replyError(http.StatusNotFound, w, r, "Not Found")
		return
	}

	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		replyError(http.StatusBadRequest, w, r, err.Error())
		return
	}
	defer r.Body.Close()

	usersIndex[user.UserName] = &user

	replyOk(w)
}

func UserDel(w http.ResponseWriter, r *http.Request) {
	userName := mux.Vars(r)["user"]

	if (!userExists(userName)) {
		replyError(http.StatusBadRequest, w, r, "Bad Request")
		return
	}

	delete(usersIndex, userName)

	replyOk(w)
}

func UserSearch(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		replyError(http.StatusBadRequest, w, r, "Bad Request")
		return
	}

	var response []*User
	for _, user := range usersIndex {
		if (strings.Compare(country, user.Country) == 0) {
			response = append(response, user)
		}
	}

	replyJSON(w, http.StatusOK, response)
}

func userExists(user string) bool {
	if _, result := usersIndex[user]; result {
		return true
	}
	return false
}

func replyJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func replyError(err int, w http.ResponseWriter, r *http.Request, text string) {
	replyJSON(w, err, map[string]string{"error": text})
}

func replyOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type",
		"application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
