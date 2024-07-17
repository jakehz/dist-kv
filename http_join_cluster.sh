#! /bin/bash


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"ipAddress":"node2:1235"}' \
  http://localhost:8080/join
