package main

import (
	"chatapp/roommanager"
	"chatapp/ws"
	"chatapp/pubsubmanager"
	
)

func main(){
	rm := roommanager.CreateRoomManagerSingleton()
	pubsub := pubsubmanager.CreatePubSubInstance(rm) // pubsub required the brodcast of rm dependency inecjection 
	// made them into a global contracts 
	rm.Pub = pubsub
    rm.Sub = pubsub
	go rm.StartRoomManager()
	go ws.CreateServer()
	
	select {} // block forever

}