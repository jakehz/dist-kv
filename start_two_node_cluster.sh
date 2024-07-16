#! /bin/bash

go run . node1 localhost 1234 8080 & go run . node2 localhost 1235 8081 && fg
