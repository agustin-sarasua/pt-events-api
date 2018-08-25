cd lambda
env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/main.zip /tmp/main

aws lambda create-function \
  --region us-east-1 \
  --function-name events-api \
  --memory 128 \
  --role arn:aws:iam::161262005667:role/EventsAPILambdaRole \
  --runtime go1.x \
  --zip-file fileb:///tmp/main.zip \
  --handler main