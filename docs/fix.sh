#!/bin/sh

sed -i '' -e "s/x-nullable/nullable/g" ./docs/docs.go
sed -i '' -e "s/x-omitempty/omitempty/g" ./docs/docs.go
sed -i '' -e "s/x-nullable/nullable/g" ./docs/swagger.json
sed -i '' -e "s/x-omitempty/omitempty/g" ./docs/swagger.json
sed -i '' -e "s/x-nullable/nullable/g" ./docs/swagger.yaml
sed -i '' -e "s/x-omitempty/omitempty/g" ./docs/swagger.yaml

sed -i '' -e "s/\"null\"/null/g" ./docs/docs.go
sed -i '' -e "s/\"null\"/null/g" ./docs/swagger.json
sed -i '' -e "s/\"null\"/null/g" ./docs/swagger.yaml
