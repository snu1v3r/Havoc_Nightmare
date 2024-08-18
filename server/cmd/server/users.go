package server

import (
	"Havoc/pkg/logger"
	"Havoc/pkg/profile"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
)

const (
	UserStatusOffline = 0
	UserStatusOnline  = 1
)

func (t *Teamserver) UserAuthenticate(username, password string) bool {
	var (
		operators []profile.Operator
		err       error
	)

	operators, err = t.profile.Operators()
	if err != nil {
		logger.DebugError("failed to get operators from profile: %v", err)
		return false
	}

	for _, operator := range operators {
		if operator.Username == username && operator.Password == password {
			return true
		}
	}

	return false
}

func (t *Teamserver) UserLogin(token string, login any, socket *websocket.Conn) {
	var (
		user string
		err  error
	)

	// get the username of the login data
	// sanity check if it's the right types
	switch login.(type) {
	case map[string]any:

		switch login.(map[string]any)["username"].(type) {
		case string:
			user = login.(map[string]any)["username"].(string)
			break
		default:
			logger.DebugError("UserLogin: invalid login username type")
			return
		}

		break

	default:
		logger.DebugError("UserLogin: invalid login type")
		return
	}

	// store the token and client data & socket into the sync map
	t.clients.Store(token, &HavocUser{
		socket:   socket,
		username: user,
	})

	// retrieve all saved events
	// and send them to the logged-in user
	t.events.mutex.Lock()
	for _, event := range t.events.list {
		if err = t.UserSend(user, event); err != nil {
			logger.DebugError("failed to send user event: %v", err)
			return
		}
	}
	t.events.mutex.Unlock()

	// broadcast that we have new client connected
	t.UserBroadcast(false, t.EventCreate(EventUserLogin, map[string]any{
		"username": user,
	}))
}

func (t *Teamserver) UserLogoutByToken(token string) error {
	var (
		value any
		ok    bool
		user  *HavocUser
		err   error
	)

	if value, ok = t.clients.LoadAndDelete(token); !ok {
		return errors.New("user was not found by token")
	}

	user = value.(*HavocUser)

	// close the socket connection in
	// case it's not closed already
	err = user.socket.Close()

	// broadcast to every connected client
	// that this user has been disconnected
	t.UserBroadcast(false, t.EventCreate(EventUserLogout, map[string]any{
		"username": user.username,
	}))

	return err
}

func (t *Teamserver) UserStatus(username string) int {
	var status = UserStatusOffline

	t.clients.Range(func(key, value any) bool {
		if value.(*HavocUser).username == username {
			status = UserStatusOnline
			return false
		}
		return true
	})

	return status
}

func (t *Teamserver) UserNameByToken(token string) (string, error) {
	var (
		name string
		err  error
	)

	err = errors.New("user not found by token")

	t.clients.Range(func(key, value any) bool {
		var user = value.(*HavocUser)
		if key == token {
			name = user.username
			err = nil
			return false
		}

		return true
	})

	return name, err
}

func (t *Teamserver) UserSend(username string, event map[string]any) error {
	var (
		user *HavocUser
		err  error
		data []byte
	)

	// default message
	err = errors.New("user not found")

	// iterate over connected users/clients
	t.clients.Range(func(key, value any) bool {
		user = value.(*HavocUser)

		// check if the connected client user
		// is the one we are searching for
		if user.username == username {
			// marshal the event map into a json byte array
			if data, err = json.Marshal(event); err != nil {
				return false
			}

			user.mutex.Lock()
			defer user.mutex.Unlock()

			// write the message to the specified user
			if err = user.socket.WriteMessage(websocket.BinaryMessage, data); err != nil {
				return false
			}

			// no error received so overwrite the one
			// we specified at the beginning of the function
			err = nil

			return false
		}

		return true
	})

	return err
}

func (t *Teamserver) UserBroadcast(save bool, event map[string]any) {
	var (
		user *HavocUser
		err  error
		data []byte
	)

	if event == nil {
		return
	}

	// iterate over connected users/clients
	t.clients.Range(func(key, value any) bool {
		user = value.(*HavocUser)

		// marshal the event map into a json byte array
		if data, err = json.Marshal(event); err != nil {
			return false
		}

		user.mutex.Lock()
		defer user.mutex.Unlock()

		// write the message to the user
		if err = user.socket.WriteMessage(websocket.BinaryMessage, data); err != nil {
			return false
		}

		return true
	})

	// if the caller specified to save the event
	// then append it to the event list
	if save {
		t.events.mutex.Lock()
		t.events.list = append(t.events.list, event)
		t.events.mutex.Unlock()
	}
}
