#! /bin/bash

export KEY="test"
export VAL="testval"

curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/get_memberlist


