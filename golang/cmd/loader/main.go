package main

import (
	"flag"
	"log"

	"github.com/mikan/mqtt-benchmarks/golang"
)

var broker = flag.String("broker", "localhost:1883", "MQTT broker url")

func main() {
	flag.Parse()
	log.SetPrefix("[LOADER] ")
	client := bench.NewClient(*broker)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	log.Print(client.ListenAndLoad())
}
