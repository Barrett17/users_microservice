package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

var url = "http://127.0.0.1:8080/api/users/"

func TestMain(t *testing.T) {
	user := User{
		"pippo",
		"pluto",
		"pippo.pluto",
		"strongpass",
		"someemail@email.com",
		"Italy"}

	addUserAndTest(t, user)

	testStressAddUsers(t)
	testStressUpdateUsers(t)
	testGetUsers(t)
	testStressDelUsers(t)

	delUserAndTest(t, user)
}

func testStressAddUsers(t *testing.T) {
	for i := 0; i < 100; i++ {
		addUserAndTest(t,
			User{
				"thx",
				"1138",
				strconv.Itoa(i),
				"strongpass",
				"someemail@email.com",
				"Italy"})
	}
}

func testStressUpdateUsers(t *testing.T) {
	for i := 0; i < 100; i++ {
		updateUserAndTest(t,
			User{
				"kyuss",
				"DnD",
				// We use the index as nickname to ease stress testing
				strconv.Itoa(i),
				"strongpass",
				"someemail@email.com",
				"UK"})
	}
}

func testStressDelUsers(t *testing.T) {
	for i := 0; i < 100; i++ {
		delUserAndTest(t,
			User{
				"thx",
				"1138",
				strconv.Itoa(i),
				"strongpass",
				"someemail@email.com",
				"Italy"})
	}
}

func testGetUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", url+"/search?country=Italy", nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Cannot GET user!")
	}

	var response []*User

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf("Cannot decode user list!")
	}

	if response == nil {
		t.Errorf("Response should not be null!")
	}

	if len(response) != 1 {
		t.Errorf("Expected one user in the list!")
	}
}

func addUserAndTest(t *testing.T, user User) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(user)
	req, _ := http.NewRequest("POST", url, buf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Cannot add user!")
	}
}

func delUserAndTest(t *testing.T, user User) {
	req, _ := http.NewRequest("DELETE", url+user.UserName, nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Cannot del user!")
	}
}

func updateUserAndTest(t *testing.T, user User) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(user)
	req, _ := http.NewRequest("PUT", url+user.UserName, buf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Cannot update user!")
	}
}
