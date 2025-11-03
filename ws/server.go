package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"fmt"
	"encoding/json"
	"chatapp/contracts"
)
var upgrader = websocket.Upgrader{}
// way of doing enums here 

type ClientMessage struct {
	Socket *websocket.Conn
	Payload contracts.Message
}

var MessageChanel = make(chan ClientMessage , 100);
func wsHandler(c *gin.Context){
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil);
	if err != nil {
		return 
	}
	defer conn.Close() ;
	var m contracts.Message ;
	
	for {
		
		_ , p , err := conn.ReadMessage() 
		if err!=nil{
			fmt.Print("error reading messages")
		}
		if err := json.Unmarshal(p, &m); err != nil {
			fmt.Println("json error:", err)
			continue
		}
		//fmt.Println("recived message" );
		//fmt.Print(m)
		MessageChanel <- ClientMessage{Socket :conn , Payload : m} ; 
		//fmt.Print(MessageChanel);
	}
	
}

func CreateServer() {
  router := gin.Default()
  router.GET("/ws" , func(c *gin.Context) {
	wsHandler(c);
  })
  router.Run() // listens on 0.0.0.0:8080 by default
}

