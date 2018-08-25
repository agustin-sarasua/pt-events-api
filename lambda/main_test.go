package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

// Define a mock struct to be used in your unit tests
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, nil
}

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

// func (m *mockDynamoDBClient) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
// 	// mock response/functionality
// }

func TestCreate(t *testing.T) {
	// Setup Test
	db = &mockDynamoDBClient{}

	auth := map[string]interface{}{
		"claims": map[string]interface{}{
			"sub": "example-sub",
		},
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  events.APIGatewayProxyResponse
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{
				Body:           "{\"name\":\"Salon de Case 2\"}",
				HTTPMethod:     "POST",
				Path:           "/persons",
				RequestContext: events.APIGatewayProxyRequestContext{Authorizer: auth},
			},
			expect: events.APIGatewayProxyResponse{StatusCode: 202, Body: "{\"name\":\"Salon de Case 2\",\"sub\":\"example-sub\"}"},
			err:    nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect.StatusCode, response.StatusCode)
		assert.Equal(t, test.expect.Body, response.Body)
	}
}
