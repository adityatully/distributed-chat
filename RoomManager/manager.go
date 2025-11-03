package roommanager

import (
	"chatapp/contracts"
	"chatapp/ws"
	"fmt"
	"sync"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var RoomManagerInstance *RoomManager // a singelton instance 
var Once sync.Once 
type Client struct {
    Conn    *websocket.Conn
    writeLock sync.Mutex   // this mutex to enure that if different go routines are writing to the same connection object 
	// then ek ek krke krene // when two diff people message in the same group at the same time 
}

func (c *Client) WriteMessage(messageType int, data []byte) error {
    c.writeLock.Lock()
    defer c.writeLock.Unlock()
    return c.Conn.WriteMessage(messageType, data)
}

type RoomManager struct {
	rooms map[string][]*Client // connection objects are ptrs 
	mu sync.RWMutex  // this because their can be read write race conditions while updating the map becuase 
	// it is hapoening in differnet go routines 
	Pub contracts.Publisher
	Sub contracts.Subscriber
}

func CreateRoomManagerSingleton() *RoomManager{
	Once.Do(func ()  {
		RoomManagerInstance = &RoomManager{
			rooms : make(map[string][]*Client) , 
			Pub : nil , 
			Sub : nil ,
		}
	})
	return RoomManagerInstance;
}


func GetRoomManagerInstance() *RoomManager{
	return RoomManagerInstance ;
}

func (rm *RoomManager)StartRoomManager(){
	for message := range ws.MessageChanel {
		msg := message.Payload 
		fmt.Println(message);
		if message.Payload.Type == "create"{
			rm.CreateRoom(msg.RoomID , message.Socket);
		}else if message.Payload.Type == "join"{
			rm.HandleJoin(msg , message.Socket)
		}else {
			rm.HandleSend(msg)
		}
	}
}

func (rm *RoomManager) CreateRoom(RoomID string , conn *websocket.Conn){
	defer rm.mu.Unlock()
	rm.mu.Lock()
	clients := []*Client{}
	if conn != nil {
		clients = append(clients, &Client{Conn: conn})
	}
	rm.rooms[RoomID] = clients
}

func(rm *RoomManager) HandleJoin(msg contracts.Message , conn *websocket.Conn){
	rm.mu.RLock() 
	_ , ok := rm.rooms[msg.RoomID]; 
	rm.mu.RUnlock()

	if !ok{
		// value dosen exist , we create a local rooms 
		rm.CreateRoom(msg.RoomID , conn);
		// we need to subscibe to this room on the PubSUB 
		go rm.Sub.SubscribeToRoom(msg.RoomID);
		return 
	}

	rm.mu.Lock()
	rm.rooms[msg.RoomID] = append(rm.rooms[msg.RoomID], &Client{Conn: conn})  // add to the room
	rm.mu.Unlock()
}

func (rm *RoomManager) HandleSend(msg contracts.Message){
	// sedn message will always come after join so a local room always exists 
	data, _ := json.Marshal(msg)
	rm.BroadcastLocal(msg.RoomID, data)

	if rm.Pub != nil {
		_ = rm.Pub.PublishToRoom(msg.RoomID, data)
	}
}

func(rm *RoomManager) BroadcastLocal(RoomID string , data []byte){
	rm.mu.RLock()
	clients := append([]*Client(nil) , rm.rooms[RoomID]...)  // copied the slice insted of keepig 
	// the read lock because brodcasting can take time 
	rm.mu.RUnlock()
	for _, client := range clients {
		go client.WriteMessage(websocket.TextMessage, data)
	}
}

func (rm *RoomManager) BroadCastFromRemote(msg contracts.Message) {
	data, _ := json.Marshal(msg)
	rm.BroadcastLocal(msg.RoomID, data)
}


