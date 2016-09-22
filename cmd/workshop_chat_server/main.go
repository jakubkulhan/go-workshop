package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jakubkulhan/go-workshop/chat"
)

var users map[int]*chat.User = map[int]*chat.User{}
var usersMutex sync.RWMutex

type idPair struct {
	FirstUserID  int
	SecondUserID int
}

var conversations map[idPair]*chat.Conversation = map[idPair]*chat.Conversation{}
var conversationsMutex sync.RWMutex

func main() {
	http.HandleFunc("/me", postMe)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/chat", postChat)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}
}

func postMe(w http.ResponseWriter, r *http.Request) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	me := &chat.User{}
	if err := chat.ReadRequest(r, me); err != nil {
		chat.WriteErrorResponse(w, r, chat.ErrBadRequest, err)
		return
	}

	if me.Name == "" {
		chat.WriteErrorResponse(w, r, chat.ErrUserNameMissing, errors.New("name missing"))
		return
	}

	for _, user := range users {
		if user.Name == me.Name {
			chat.WriteErrorResponse(w, r, chat.ErrUserNameAlreadyExists, fmt.Errorf("user [%s] already exists", me.Name))
			return
		}
	}

	me.ID = len(users) + 1
	users[me.ID] = me

	chat.WriteOkResponse(w, r, &chat.MeResponse{true, me})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	userList := []*chat.User{}
	for _, user := range users {
		userList = append(userList, user)
	}

	chat.WriteOkResponse(w, r, &chat.UserListResponse{true, userList})
}

func postChat(w http.ResponseWriter, r *http.Request) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	conversationsMutex.Lock()
	defer conversationsMutex.Unlock()

	message := &chat.Message{}
	if err := chat.ReadRequest(r, message); err != nil {
		chat.WriteErrorResponse(w, r, chat.ErrBadRequest, err)
		return
	}

	if _, ok := users[message.FromUserID]; !ok {
		chat.WriteErrorResponse(w, r, chat.ErrFromUserIDMissing, errors.New("unknown from user ID"))
		return
	}

	if _, ok := users[message.ToUserID]; !ok {
		chat.WriteErrorResponse(w, r, chat.ErrToUserIDMissing, errors.New("unknown to user ID"))
		return
	}

	if message.Text == "" {
		chat.WriteErrorResponse(w, r, chat.ErrTextMissing, errors.New("message text missing"))
		return
	}

	if message.FromUserID == message.ToUserID {
		chat.WriteErrorResponse(w, r, chat.ErrCannotChat, errors.New("you cannot chat with yourself"))
		return
	}

	conversationId := idPair{}
	if message.FromUserID < message.ToUserID {
		conversationId.FirstUserID = message.FromUserID
		conversationId.SecondUserID = message.ToUserID
	} else {
		conversationId.FirstUserID = message.ToUserID
		conversationId.SecondUserID = message.FromUserID
	}

	message.CreatedAt = time.Now()

	conversation, ok := conversations[conversationId]
	if !ok {
		conversation = &chat.Conversation{}
		conversations[conversationId] = conversation
	}
	conversation.Messages = append(conversation.Messages, message)

	fmt.Printf("sent message from [%d] to [%d]: %s\n", message.FromUserID, message.ToUserID, message.Text)

	chat.WriteOkResponse(w, r, &chat.OkResponse{true})
}
