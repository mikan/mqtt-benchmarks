package com.github.mikan.bench;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.eclipse.paho.client.mqttv3.*;

import java.util.UUID;
import java.util.concurrent.CountDownLatch;
import java.util.logging.Logger;

/**
 * Client provides MQTT client operations.
 */
class Client {
    private static final String CMD_TOPIC = "bench/cmd";
    private static final String LOAD_TOPIC = "bench/load";
    private static final Logger LOG = Logger.getLogger("[BENCH]");
    private final MqttClient client;

    Client(String broker) {
        try {
            client = new MqttClient("tcp://" + broker, UUID.randomUUID().toString());
        } catch (MqttException e) {
            throw new RuntimeException("Failed to create MQTT client.", e);
        }
        var options = new MqttConnectOptions();
        try {
            client.connect(options);
        } catch (MqttException e) {
            throw new RuntimeException("Failed to connect broker.", e);
        }
    }

    /**
     * Sends start message to loader and receives traffic.
     *
     * @param nPublish number of publishes
     */
    void Bench(int nPublish) {
        var countDownLatch = new CountDownLatch(nPublish);
        try {
            client.subscribeWithResponse(LOAD_TOPIC, 1, (topic, message) -> {
                countDownLatch.countDown();
            }).waitForCompletion();
            LOG.info("subscribe " + LOAD_TOPIC);
        } catch (MqttException e) {
            throw new RuntimeException("Failed to subscribe " + LOAD_TOPIC + ".", e);
        }
        String payload;
        try {
            var msg = new CommandMessage();
            msg.nPublish = nPublish;
            payload = new ObjectMapper().writeValueAsString(msg);
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Failed to build command message.", e);
        }
        try {
            client.publish(CMD_TOPIC, payload.getBytes(), 1, false);
            LOG.info("loading...");
        } catch (MqttException e) {
            throw new RuntimeException("Failed to publish " + CMD_TOPIC, e);
        }
        try {
            countDownLatch.await();
            LOG.info("done. " + nPublish + " messages received.");
        } catch (InterruptedException e) {
            LOG.info("Interrupted: " + e.getMessage());
        }
        try {
            client.unsubscribe(LOAD_TOPIC);
            client.disconnect();
            client.close();
        } catch (MqttException e) {
            throw new RuntimeException("Failed to close connection.", e);
        }
        if (countDownLatch.getCount() != 0) {
            LOG.severe("Benchmark failed. Request " + nPublish + ", got " + (nPublish - countDownLatch.getCount()) + ".");
        }
    }

    private static class CommandMessage {
        @JsonProperty("n_publish")
        private int nPublish;
    }
}
