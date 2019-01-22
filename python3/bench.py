#!/usr/bin/env python3
import argparse
import time
import uuid

import paho.mqtt.client as mqtt


class Client:
    """
    Provides MQTT client operations.
    """

    CMD_TOPIC = "bench/cmd"
    LOAD_TOPIC = "bench/load"

    def __init__(self, host, port):
        """
        Initialize the MQTT client.

        :param host: broker host name
        :param port: broker port number
        """
        self.client = mqtt.Client(client_id=str(uuid.uuid4()), protocol=mqtt.MQTTv311)
        self.host = host
        self.port = port
        self.n_publish = 0
        self.n_receives = 0
        self.done = False

    def bench(self, n_publish):
        self.n_publish = n_publish
        self.client.on_subscribe = self.on_subscribe
        self.client.on_message = self.on_message
        self.client.on_disconnect = self.on_disconnect
        self.client.connect(self.host, port=self.port)
        self.client.subscribe(self.LOAD_TOPIC, 1)
        self.client.loop_start()
        while self.n_receives < self.n_publish:
            time.sleep(0.1)
            pass
        self.client.unsubscribe(self.LOAD_TOPIC)
        self.client.loop_stop()
        print("[BENCH] done. %d messages received." % self.n_receives)
        if self.n_receives < n_publish:
            print("[BENCH] benchmark failed. request %d, got %d.", n_publish, self.n_receives)

    def on_subscribe(self, client, _userdata, _mid, _granted_qos):
        print("[BENCH] subscribe " + self.LOAD_TOPIC)
        client.publish(self.CMD_TOPIC, """{"n_publish":%d}""" % self.n_publish)
        print("[BENCH] loading...")

    def on_message(self, _client, _userdata, _msg):
        self.n_receives = self.n_receives + 1
        if self.n_receives >= self.n_publish:
            self.done = True

    def on_disconnect(self, _client, _userdata, _rc):
        print("[BENCH] disconnected from %s:%d" % (self.host, self.port))


def main():
    parser = argparse.ArgumentParser(description="MQTT benchmark program.")
    parser.add_argument("-host", type=str, default="localhost", help="MQTT broker host name")
    parser.add_argument("-port", type=int, default=1883, help="MQTT broker port number")
    parser.add_argument("-n", type=int, default=1000, help="Number of publishes")
    args = parser.parse_args()
    client = Client(args.host, args.port)
    client.bench(args.n)


if __name__ == '__main__':
    main()
