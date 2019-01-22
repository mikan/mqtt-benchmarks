#!/bin/bash

function usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "command:"
    echo "  build   Build all programs"
    echo "  loader  Run loader program"
    echo "  g <n>   Run golang bench"
    echo "  j <n>   Run java bench"
    echo "  p <n>   Run python bench"
}

if [[ $# == 0 ]]; then
    usage
fi

if [[ $# == 1 ]]; then
    if [[ $1 == "build" ]]; then
        echo "Building golang program.., (1 of 2)"
        cd golang
        make build
        echo "Building java program... (2 of 2)"
        cd ../java11
        ./gradlew shadowJar
        cd ..
    elif [[ $1 == "loader" ]]; then
        ./golang/build/loader &
    elif [[ $1 == "g" ]]; then
        ./golang/build/bench
    elif [[ $1 == "j" ]]; then
        java -jar java11/build/libs/mqtt-benchmark.jar
    elif [[ $1 == "p" ]]; then
        ./python3/bench.py
    else
        usage
    fi
elif [[ $# == 2 ]]; then
    if [[ $1 == "g" ]]; then
        ./golang/build/bench -n $2
    elif [[ $1 == "j" ]]; then
        java -jar java11/build/libs/mqtt-benchmark.jar -n $2
    elif [[ $1 == "p" ]]; then
        ./python3/bench.py -n $2
    else
        usage
    fi
else
    usage
fi
