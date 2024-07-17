#! /bin/bash

export KEY="test"
export VAL="testval"

curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8080/set/$KEY/$VAL
echo "\n"
echo "Getting 8080"
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/get_whole_kv

echo "\n"
echo "Getting 8081"
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8081/get_whole_kv

