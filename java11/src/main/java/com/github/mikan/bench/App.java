package com.github.mikan.bench;

import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.Options;
import org.apache.commons.cli.ParseException;

/**
 * MQTT benchmark program.
 */
public class App {
    public static void main(String[] args) {
        var broker = "localhost:1883";
        var nPublish = 1000;
        var options = new Options();
        options.addOption("broker", true, "MQTT broker url");
        options.addOption("n", true, "Number of publishes");
        var parser = new DefaultParser();
        try {
            var cl = parser.parse(options, args);
            broker = cl.getOptionValue("broker", broker);
            var nPublishStr = cl.getOptionValue("n", Integer.toString(nPublish));
            try {
                nPublish = Integer.parseInt(nPublishStr);
            } catch (NumberFormatException e) {
                System.err.println("invalid n: " + nPublishStr);
                System.exit(2);
            }
        } catch (ParseException e) {
            System.err.println("invalid parameter: " + e.getMessage());
            System.exit(2);
        }
        var client = new Client(broker);
        client.Bench(nPublish);
    }
}
