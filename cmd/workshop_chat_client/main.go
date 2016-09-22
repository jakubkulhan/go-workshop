package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jakubkulhan/go-workshop/chat"
)

var commands map[string]func(string) = map[string]func(string){
	"/list": list,
	"/to":   to,
}
var myUserID int = 0
var otherUserID int = 0

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Your name?")
	fmt.Print("> ")

	if !scanner.Scan() {
		return
	}

	buf, err := json.Marshal(&chat.User{Name: scanner.Text()})
	if err != nil {
		panic(err)
	}
	response, err := http.Post("http://localhost:8080/me", "application/json", bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	meResponse := &chat.MeResponse{}
	if err := json.Unmarshal(mustReadAll(response.Body), meResponse); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", meResponse)
	myUserID = meResponse.Data.ID

	// TODO: check meResponse.OK

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		if strings.HasPrefix(scanner.Text(), "/") {
			parts := strings.SplitN(scanner.Text(), " ", 2)

			handler, ok := commands[parts[0]]
			if !ok {
				fmt.Printf("command [%s] not found\n", parts[0])
			} else if len(parts) > 1 {
				handler(parts[1])
			} else {
				handler("")
			}

		} else if otherUserID != 0 {
			buf, err := json.Marshal(&chat.Message{
				FromUserID: myUserID,
				ToUserID:   otherUserID,
				Text:       scanner.Text(),
			})
			if err != nil {
				panic(err)
			}
			response, err := http.Post("http://localhost:8080/chat", "application/json", bytes.NewReader(buf))
			if err != nil {
				panic(err)
			}
			chatResponse := &chat.ErrorResponse{}
			if err := json.Unmarshal(mustReadAll(response.Body), chatResponse); err != nil {
				panic(err)
			}

			if !chatResponse.OK {
				fmt.Printf("error: %s\n", chatResponse.Message)
			}

		} else {
			fmt.Println("no open conversation")
		}
	}
}

func list(line string) {
	response, err := http.Get("http://localhost:8080/users")
	if err != nil {
		panic(err)
	}
	userListResponse := &chat.UserListResponse{}
	if err := json.Unmarshal(mustReadAll(response.Body), userListResponse); err != nil {
		panic(err)
	}

	// TODO: check userListResponse.OK

	for _, user := range userListResponse.Data {
		fmt.Printf("[%d] %s\n", user.ID, user.Name)
	}
}

func to(line string) {
	id, err := strconv.ParseInt(strings.TrimSpace(line), 10, 32)
	if err != nil {
		panic(err)
	}
	otherUserID = int(id)
}

func mustReadAll(r io.Reader) []byte {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return buf
}
