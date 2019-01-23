package main

import (
	"flag"
	"log"

	"github.com/mikan/mqtt-benchmarks/golang"
)

var broker = flag.String("broker", "localhost:1883", "MQTT broker url")
var nPublish = flag.Int("n", 1000, "Number of publishes")
var gapMilliseconds = flag.Int("gap", 0, "Gap time between publishes as milliseconds")

func main() {
	flag.Parse()
	log.SetPrefix("[BENCH] ")
	client := bench.NewClient(*broker)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	if err := client.Bench(*nPublish, *gapMilliseconds); err != nil {
		log.Fatal(err)
	}
}
