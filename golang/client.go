package bench

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const (
	cmdTopic  = "bench/cmd"
	loadTopic = "bench/load"
)

// Client implements MQTT client.
type Client struct {
	client        mqtt.Client
	nPublish      int
	nReceives     int
	nReceivesLock sync.Mutex
	done          chan struct{}
}

type cmdMessage struct {
	NPublish int `json:"n_publish"`
}

// NewClient will creates MQTT client instance.
func NewClient(broker string) *Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + broker)
	opts.SetClientID(uuid.New().String())
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Fatalf("connection lost: %v", err)
	})
	client := &Client{
		client: mqtt.NewClient(opts),
	}
	return client
}

// Connect opens MQTT broker connection.
func (c *Client) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect broker: %v", token.Error())
	}
	return nil
}

// Disconnect closes MQTT broker connection.
func (c *Client) Disconnect() {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to disconnect broker: %v", token.Error())
	}
}

// ListenAndLoad listens command topic and sends traffic.
func (c *Client) ListenAndLoad() error {
	if token := c.client.Subscribe(cmdTopic, 1, handleLoad); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe %s: %v", cmdTopic, token.Error())
	}
	log.Printf("subscribe %s", cmdTopic)
	for {
		time.Sleep(1 * time.Second)
	}
}

func handleLoad(client mqtt.Client, msg mqtt.Message) {
	var cmd cmdMessage
	if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
		log.Printf("failed to parse cmd message: %v", err)
	}
	log.Printf("%s received. n=%d", cmdTopic, cmd.NPublish)
	for i := 0; i < cmd.NPublish; i++ {
		time.Sleep(1 * time.Millisecond)
		payload := fmt.Sprintf(`{"count":%d}`, i+1)
		if token := client.Publish(loadTopic, 1, false, payload); token.Wait() && token.Error() != nil {
			log.Printf("failed to publish %s: %v", loadTopic, token.Error())
		}
	}
}

// Bench sends message to loader and receives traffic.
func (c *Client) Bench(nPublish int) error {
	c.done = make(chan struct{})
	c.nPublish = nPublish
	if token := c.client.Subscribe(loadTopic, 1, c.handleBench); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe %s: %v", loadTopic, token.Error())
	}
	log.Printf("subscribe %s", loadTopic)
	payload, err := json.Marshal(cmdMessage{NPublish: nPublish})
	if err != nil {
		return fmt.Errorf("failed to build message: %v", err)
	}
	if token := c.client.Publish(cmdTopic, 1, false, payload); token.Wait() && token.Error() != nil {
		log.Printf("failed to publish %s: %v", cmdTopic, token.Error())
	}
	log.Print("loading...")
	<-c.done
	c.nReceivesLock.Lock()
	defer c.nReceivesLock.Unlock()
	if token := c.client.Unsubscribe(loadTopic); token.Wait() && token.Error() != nil {
		log.Printf("failed to unsubscribe %s: %v", loadTopic, token.Error())
	}
	log.Printf("done. %d messages received.", c.nReceives)
	if c.nReceives < nPublish {
		log.Printf("benchmark failed. request %d, got %d.", nPublish, c.nReceives)
	}
	return nil
}

func (c *Client) handleBench(_ mqtt.Client, msg mqtt.Message) {
	c.nReceivesLock.Lock()
	c.nReceives++
	c.nReceivesLock.Unlock()
	if c.nReceives >= c.nPublish {
		c.done <- struct{}{}
	}
}
