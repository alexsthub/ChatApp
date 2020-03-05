package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/user"
	"sync"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// Notifier is a struct to save websocket connections
type Notifier struct {
	mu          sync.RWMutex
	connections map[int64]*websocket.Conn
}

// NewNotifier ah
func NewNotifier() *Notifier {
	ws := &Notifier{
		mu:          sync.RWMutex{},
		connections: make(map[int64]*websocket.Conn),
	}
	return ws
}

func (notifier *Notifier) addConnection(userID int64, conn *websocket.Conn) {
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	notifier.connections[userID] = conn
}

func (notifier *Notifier) getConnection(userID int64) *websocket.Conn {
	notifier.mu.RLock()
	defer notifier.mu.RUnlock()
	conn, exists := notifier.connections[userID]
	if !exists {
		return nil
	}
	return conn
}

func (notifier *Notifier) removeConnection(userID int64) {
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	delete(notifier.connections, userID)
}

// WebSocketConnectionHandler upgrade the connection to a web socket connection if the user is authenticated
func (ctx *ContextHandler) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, "User not authenticated: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Header.Get("Origin") != "https://alexst.me" {
				return false
			}
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to open websocket connection", 401)
		return
	}
	ctx.Notifier.addConnection(sessionState.User.ID, conn)
}

type channel struct {
	_id         string
	name        string
	description string
	private     bool
	members     []*user.User
	createdAt   time.Time
	creator     *user.User
	editedAt    time.Time
}

type message struct {
	_id       string
	channelID string
	body      string
	createdAt time.Time
	creator   *user.User
	editedAt  time.Time
}

type messageObj struct {
	messageType string
	channel     *channel
	channelID   string
	message     *message
	messageID   string
	userIDs     []int64
}

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

// RabbitMQConn shit
func (ctx *ContextHandler) RabbitMQConn() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Print("Failed to connecto rabbitMQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Print("Failed to create channel")
	}

	q, err := ch.QueueDeclare(
		"message", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Print("Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go ctx.Notifier.run(msgs)
}

func (notifier *Notifier) run(msgs <-chan amqp.Delivery) {
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	for msg := range msgs {
		message := &messageObj{}
		err := json.Unmarshal([]byte(msg.Body), message)
		if err != nil {
			fmt.Println(err)
		}
		if len(message.userIDs) > 0 {
			// If private
			for _, id := range message.userIDs {
				socket := notifier.getConnection(id)
				err = socket.WriteMessage(websocket.TextMessage, msg.Body)
				if err != nil {
					notifier.removeConnection(id)
				}
			}
		} else {
			// Public
			for id, socket := range notifier.connections {
				err = socket.WriteMessage(websocket.TextMessage, msg.Body)
				if err != nil {
					notifier.removeConnection(id)
				}
			}
		}
	}
}
