cd lambda
env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/main.zip /tmp/main

# --handler is the path to the executable inside the .zip
aws lambda update-function-code --function-name events-api \
--zip-file fileb:///tmp/main.zip

# Test the deploy
# > aws lambda invoke --function-name CreateEvent /tmp/output.json
# > cat /tmp/output.json