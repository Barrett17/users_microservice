package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	USER_ADDED = iota
	USER_CHANGED_NICK = iota
)

type UserExternalNotification struct {
	what int `json:"what"`
	data string `json:"data"`
}

// At this point I expect to have a notification server (custom or 
// built around kafka or other kind of message queues) if we have
// the need to broadcast the message in a 1:n fashion.
// Otherwise, like in this case it seems we have a 1:1 configuration
// with other microservices, so I expect other servers to add their own
// REST endpoints which we just invoke to notify them of changes.
// I could have used the OS IPC for example, but this would not have
// allowed to easily scale up on different machines. That's why this
// is in this case the only reasonable solution to the notification
// problems. Obviously this code will fail, that's why I don't add checks
// on the caller side, to avoid this test failing for no reasons.

var searchService = "https://thesearchservice/api/search"
var competitionService = "https://thecompetitionservice.com/api/competitions"

func UserChanged(user User) {
	message := UserExternalNotification{USER_CHANGED_NICK, user.UserName}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(message)
	req, _ := http.NewRequest("PUT", competitionService, buf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
}

func UserAdded(user User) {
	message := UserExternalNotification{USER_ADDED, user.UserName}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(message)
	req, _ := http.NewRequest("POST", searchService, buf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
}
