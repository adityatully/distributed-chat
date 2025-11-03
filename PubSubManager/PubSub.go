package pubsubmanager

import (
	"context"
	"fmt"
	"sync"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"chatapp/contracts"
)

var ready = make(chan struct{}) // channel for signaloling to all other like await 
var PubSubInstance *PubSubManager
var Once sync.Once

type PubSubManager struct {
	rclient 		*redis.Client 
	BroadCaster 		contracts.BroadCaster
	subscriptions 	map[string]struct{}
	mu 				sync.Mutex
}

func CreatePubSubInstance(brodaster contracts.BroadCaster) *PubSubManager{
	Once.Do(func(){
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		PubSubInstance = &PubSubManager{
			rclient: client,
			BroadCaster: brodaster,
			subscriptions:  make(map[string]struct{}),
		}
		close(ready)
	})
	<-ready ; 
	return PubSubInstance
}


func GetSingletonInstance()*PubSubManager{
	return PubSubInstance
}


func (ps *PubSubManager) PublishToRoom(channel string, message []byte) error {
	return ps.rclient.Publish(context.Background(), channel, message).Err()
}

func (ps *PubSubManager)SubscribeToRoom(channelName string){
	ps.mu.Lock()
	if _, already := ps.subscriptions[channelName]; already {
		ps.mu.Unlock()
		return
	}
	ps.subscriptions[channelName] = struct{}{}
	ps.mu.Unlock()

	pubsub := ps.rclient.Subscribe(context.Background(), channelName)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var m contracts.Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				fmt.Println("Redis message unmarshal error:", err)
				continue
			}
			// Notify RoomManager (via Broadcaster interface)
			ps.BroadCaster.BroadCastFromRemote(m)
		}
	}()
}