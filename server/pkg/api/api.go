package api

import (
	"Havoc/pkg/cert"
	"Havoc/pkg/colors"
	"Havoc/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	ApiMaxRequestRead = 50000000 // 50mb
	ApiTokenLength    = 32
	ApiTokenHeader    = "x-havoc-token"
)

const (
	UserStatusOffline = 0
	UserStatusOnline  = 1
)

type HavocInterface interface {
	UserAuthenticate(username, password string) bool
	UserLogin(token string, login any, socket *websocket.Conn)
	UserLogoutByToken(user string) error
	UserNameByToken(token string) (string, error)
	UserStatus(username string) int

	ListenerStart(name, protocol string, options map[string]any) error
	ListenerStop(name string) error
	ListenerRestart(name string) error
	ListenerRemove(name string) error
	ListenerEvent(protocol string, event map[string]any) (map[string]any, error)
	ListenerConfig(name string) (map[string]any, error)
	ListenerEdit(name string, config map[string]any) error
	ListenerList() []map[string]string

	AgentInitialize(uuid, plugin, status, note string, data map[string]any) error
	AgentData(uuid string) (map[string]any, error)
	AgentExist(uuid string) bool
	AgentType(uuid string) (string, error)
	AgentRemove(uuid string) error
	AgentList() []string
	AgentGenerate(implant string, config map[string]any) (string, []byte, map[string]any, error)
	AgentExecute(uuid string, data map[string]any, wait bool) (map[string]any, error)
	AgentNote(uuid string) (string, error)
	AgentSetNote(uuid string, note string) error
	AgentStatus(uuid string) (string, error)
	AgentSetStatus(uuid string, status string) error

	DatabaseAgentConsole(uuid string) ([]map[string]any, error)
}

type ServerApi struct {
	Engine *gin.Engine

	// ssl structure containing the path to
	// the generated ssl cert and key file
	ssl struct {
		cert string
		key  string
	}

	// havoc interface
	// to interact with some functions to
	// add/query/remove objects and data
	// to from the teamserver
	havoc HavocInterface

	// wait queue for websockets to send
	// the access token to register the
	// event endpoint
	tokens sync.Map
}

func NewServerApi(teamserver HavocInterface) (*ServerApi, error) {
	var api = new(ServerApi)

	gin.SetMode(gin.ReleaseMode)

	api.havoc = teamserver
	api.Engine = gin.New()

	//
	// set api endpoints
	//
	api.Engine.POST("/api/login", api.login)

	//
	// listeners endpoints
	//
	api.Engine.POST("/api/listener/list", api.listenerList)
	api.Engine.POST("/api/listener/start", api.listenerStart)
	api.Engine.POST("/api/listener/stop", api.listenerStop)
	api.Engine.POST("/api/listener/restart", api.listenerRestart)
	api.Engine.POST("/api/listener/remove", api.listenerRemove)
	api.Engine.POST("/api/listener/edit", api.listenerEdit)
	api.Engine.POST("/api/listener/event", api.listenerEvent)
	api.Engine.POST("/api/listener/config", api.listenerConfig)

	//
	// agent endpoints
	//
	api.Engine.POST("/api/agent/list", api.agentList)
	api.Engine.POST("/api/agent/build", api.agentBuild)
	api.Engine.POST("/api/agent/execute", api.agentExecute)
	api.Engine.POST("/api/agent/note", api.agentNote)
	api.Engine.POST("/api/agent/console", api.agentConsole)
	api.Engine.POST("/api/agent/remove", api.agentRemove)
	api.Engine.POST("/api/agent/hide", api.agentHide)

	//
	// websocket event endpoint
	//
	api.Engine.GET("/api/event", api.event)

	return api, nil
}

// Start the server api
// generates an HTTP certificate and then starts the api server
func (api *ServerApi) Start(host, port string) {
	var err error

	logger.Info("starting server on %v", colors.BlueUnderline("https://"+host+":"+port))

	// start the api server
	if err = api.Engine.RunTLS(host+":"+port, api.ssl.cert, api.ssl.key); err != nil {
		logger.Error("Failed to start webserver: " + err.Error())
		return
	}
}

func (api *ServerApi) GenerateSSL(host, certsPath string) (string, string, error) {
	var (
		certPath = certsPath + "/server.cert"
		keyPath  = certsPath + "/server.key"
		err      error
		Cert     []byte
		Key      []byte
	)

	// generate cert and key based on host
	Cert, Key, err = cert.HTTPSGenerateRSACertificate(host)
	if err != nil {
		logger.Error("failed to generate server certificates: " + err.Error())
		os.Exit(0)
	}

	// write cert file to disk
	err = os.WriteFile(certPath, Cert, 0644)
	if err != nil {
		return "", "", errors.New("couldn't save server cert file: " + err.Error())
	}

	// write the cert key path to disk
	err = os.WriteFile(keyPath, Key, 0644)
	if err != nil {
		return "", "", errors.New("couldn't save server cert file: " + err.Error())
	}

	return certPath, keyPath, nil
}

func (api *ServerApi) SetSSL(certPath, keyPath string) error {

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}

	api.ssl.cert = certPath
	api.ssl.key = keyPath

	return nil
}

// login endpoint for the client to log in
// and retrieve an access token
func (api *ServerApi) login(ctx *gin.Context) {
	var (
		body  []byte
		token string
		err   error
		user  string
		pass  string
		login map[string]any
	)

	// read from request the login data
	if body, err = io.ReadAll(io.LimitReader(ctx.Request.Body, ApiMaxRequestRead)); err != nil {
		logger.DebugError("Failed to read from server api login request: " + err.Error())
		goto ERROR
	}

	logger.Debug("got request on /api/login:" + fmt.Sprintf("%s", string(body)))

	// unmarshal the bytes into a map
	if err = json.Unmarshal(body, &login); err != nil {
		logger.DebugError("Failed to unmarshal bytes to map: " + err.Error())
		return
	}

	// generate a token
	// if failed then abort with status code
	if token = api.generateToken(); len(token) == 0 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get username from login request
	switch login["username"].(type) {
	case string:
		user = login["username"].(string)
	default:
		logger.DebugError("Failed retrieve username: invalid type")
		goto ERROR
	}

	// check if the specified user
	// is already active/online
	if api.havoc.UserStatus(user) == UserStatusOnline {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is already online",
		})
		return
	}

	// get password from login request
	switch login["password"].(type) {
	case string:
		pass = login["password"].(string)
	default:
		logger.DebugError("Failed retrieve password: invalid type")
		goto ERROR
	}

	// check if we can authenticate the user
	if !api.havoc.UserAuthenticate(user, pass) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// send json token response
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})

	// store the token inside the wait queue
	api.tokens.Store(token, login)

	return

ERROR:
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

// handleEventClient
// handles incoming event clients
//
// it is going to read the given token from the client and checks
// if it in the wait queue. After log in the token is going to be
// stored inside the wait queue til the client connects to the websocket
// endpoint and send the token it previously got from the /api/login
// endpoint. After all of this the teamserver is going to register
// the client as officially logged in.
//
//	socket *websocket.Conn
func (api *ServerApi) handleEventClient(socket *websocket.Conn) {
	var (
		body  []byte
		err   error
		data  map[string]any
		login any
		token string
		ok    bool
	)

	if socket == nil {
		return
	}

	// try to read the token from socket
	if _, body, err = socket.ReadMessage(); err != nil {
		logger.DebugError("failed reading message from socket: " + err.Error())
		goto ERROR
	}

	if len(body) < ApiTokenLength {
		logger.DebugError("invalid socket request. closing connection")
		goto ERROR
	}

	// unmarshal the body to a map
	if err = json.Unmarshal(body, &data); err != nil {
		logger.DebugError("failed unmarshal message from socket: " + err.Error())
		goto ERROR
	}

	// get token from request
	switch data["token"].(type) {
	case string:
		token = data["token"].(string)
	default:
		logger.DebugError("failed retrieve token: invalid type")
		goto ERROR
	}

	// check if is inside the tokens map
	if login, ok = api.tokens.Load(token); !ok {
		logger.DebugError("failed retrieve token from map: token entry not found")
		goto ERROR
	}

	// by now the client was successfully authenticated
	// we're going to register the user to the teamserver
	api.havoc.UserLogin(token, login, socket)

	// this loop is just there to check the state of the connection
	// and if the client somehow disconnected from the connection.
	// if the user disconnected then handle it by close the resources
	// related to that user session
	for {
		// check if any error arrived from reading any message or
		// checking the state of the connection
		if _, _, err = socket.ReadMessage(); err == nil {
			continue
		}

		logger.Debug("user connection closed [token: %v] -> %v", token, err)

		// log out the connected user client and send
		// a broadcast event to every connected client
		if err = api.havoc.UserLogoutByToken(token); err != nil {
			logger.DebugError("failed to logout user by token: %v", err)
		}

		break
	}

ERROR:
	socket.Close()
}

func (api *ServerApi) sanityCheck(ctx *gin.Context) bool {
	var (
		token string
		ok    bool
	)

	// check if the token header exists
	if token = ctx.GetHeader(ApiTokenHeader); len(token) == 0 || len(token) < ApiTokenLength {
		return false
	}

	// check if is inside the tokens map
	if _, ok = api.tokens.Load(token); !ok {
		logger.DebugError("failed retrieve token from map: token entry not found")
		return false
	}

	return true
}

func (api *ServerApi) tokenUser(token string) (string, bool) {
	var (
		user  string
		val   any
		login map[string]any
		ok    bool
	)

	// check if is inside the tokens map
	if val, ok = api.tokens.Load(token); !ok {
		logger.DebugError("failed retrieve token from map: token entry not found")
		return user, false
	}

	login = val.(map[string]any)

	// get name from agent request
	switch login["username"].(type) {
	case string:
		user = login["username"].(string)
	default:
		logger.DebugError("Failed retrieve login username: invalid type")
		return user, false
	}

	return user, true
}

// generateToken
// generate a token for the client to use for authenticate
// and to use to interact with the api with
func (api *ServerApi) generateToken() string {
	b := make([]byte, ApiTokenLength/2)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
