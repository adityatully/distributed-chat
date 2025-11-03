package contracts

import (
	"encoding/json"
)
type MessageType string
type Message struct {
	Type   MessageType `json:"type"`
	RoomID string      `json:"roomId"`
	Data   string      `json:"data"` 
}

func ToJson(m *Message) ([]byte , error){
	return json.Marshal(m);  // marshal makes our struct data into a json byte stream 
}

func MessageFromjson(b [] byte) (Message , error){
	var msg Message 
	err := json.Unmarshal(b , &msg)  // converts byte stream into struct 
	return msg , err
}

type Publisher interface {  // neded by the room Manager 
	// any publisher shud implement the Publish to Room method 
	PublishToRoom(channelName string , mess []byte)error
}

type Subscriber interface {  // neded by the room Manager 
	SubscribeToRoom(channelName string)
}

type BroadCaster interface { // neded by the Pubsub Manager   
	BroadCastFromRemote(msg Message)
}


