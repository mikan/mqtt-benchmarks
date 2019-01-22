mqtt benchmarks
===============

Compare MQTT client performance for Go, Java and Python.

## Build

```
cd golang
make build
cd ../java11
./gradlew shadowJar
cd ..
```

## Loader

```
./golang/build/loader &
```

## Bench

```
./golang/build/bench -n 10000
java -jar ./java11/build/libs/mqtt-benchmark.jar -n 10000
./python3/bench.py -n 10000
```
