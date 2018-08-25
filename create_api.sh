####################################################################
# Create the Rest API
####################################################################

aws apigateway create-rest-api --name 'Events API' --region us-east-1
# {
#     "apiKeySource": "HEADER",
#     "name": "Events API",
#     "createdDate": 1535116541,
#     "endpointConfiguration": {
#         "types": [
#             "EDGE"
#         ]
#     },
#     "id": "rizogxyxik"
# }

####################################################################
# Create the Resources
####################################################################

aws apigateway get-resources --rest-api-id rizogxyxik --region us-east-1
# {
#     "items": [
#         {
#             "path": "/",
#             "id": "udsvm6ycmb"
#         }
#     ]
# }

# Create resource /events
aws apigateway create-resource --rest-api-id rizogxyxik \
      --region us-east-1 \
      --parent-id udsvm6ycmb \
      --path-part events
# {
#     "path": "/events",
#     "pathPart": "events",
#     "id": "5w38av",
#     "parentId": "udsvm6ycmb"
# }

####################################################################
# Create Authorizer for the endpoints (cognito user pool)
####################################################################

# Create Authorizer (AUTH)
aws apigateway create-authorizer --rest-api-id rizogxyxik \
        --name 'Events_API_Authorizer' \
        --type COGNITO_USER_POOLS \
        --provider-arns 'arn:aws:cognito-idp:us-east-1:161262005667:userpool/us-east-1_3IwH7mpoM' \
        --identity-source 'method.request.header.Authorization' \
        --authorizer-result-ttl-in-seconds 300
# {
#     "authType": "cognito_user_pools",
#     "identitySource": "method.request.header.Authorization",
#     "name": "Events_API_Authorizer",
#     "providerARNs": [
#         "arn:aws:cognito-idp:us-east-1:161262005667:userpool/us-east-1_3IwH7mpoM"
#     ],
#     "type": "COGNITO_USER_POOLS",
#     "id": "r85iw7"
# }


# Get AUTHORIZERS
aws apigateway get-authorizers --rest-api-id rizogxyxik

####################################################################
# Create Method Request and enable responses on Resources
####################################################################

# Enable GET (AUTH) for /events
aws apigateway put-method --rest-api-id rizogxyxik \
       --resource-id 5w38av \
       --http-method GET \
       --authorization-type COGNITO_USER_POOLS \
       --authorizer-id r85iw7 \
       --region us-east-1
# {
#     "apiKeyRequired": false,
#     "httpMethod": "GET",
#     "authorizationType": "COGNITO_USER_POOLS",
#     "authorizerId": "r85iw7"
# }

# Enable POST (AUTH) for /events
aws apigateway put-method --rest-api-id rizogxyxik \
       --resource-id 5w38av \
       --http-method POST \
       --authorization-type COGNITO_USER_POOLS \
       --authorizer-id r85iw7 \
       --region us-east-1
# {
#     "apiKeyRequired": false,
#     "httpMethod": "POST",
#     "authorizationType": "COGNITO_USER_POOLS",
#     "authorizerId": "r85iw7"
# }

# Enable 200 OK to GET /events
aws apigateway put-method-response --rest-api-id rizogxyxik \
       --resource-id 5w38av --http-method GET \
       --status-code 200  --region us-east-1

# Enable 200 OK to POST /events
aws apigateway put-method-response --rest-api-id rizogxyxik \
       --resource-id 5w38av --http-method POST \
       --status-code 200  --region us-east-1

####################################################################
# Now Integrate with the Backend
# NOTE: Create the backend Lambda
####################################################################
# Create Rol for the lambda function
./create_rol.sh

# Create the lambda
./create_lambda.sh
# "FunctionArn": "arn:aws:lambda:us-east-1:161262005667:function:persons-api",

# Integrate GET /events with Lambda function
aws apigateway put-integration --rest-api-id rizogxyxik \
        --resource-id 5w38av \
        --http-method GET \
        --type AWS_PROXY \
        --integration-http-method POST \
        --region us-east-1 \
        --uri 'arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:161262005667:function:events-api/invocations'

# Integrate POST /events with Lambda function
aws apigateway put-integration --rest-api-id rizogxyxik \
        --resource-id 5w38av \
        --http-method POST \
        --type AWS_PROXY \
        --integration-http-method POST \
        --region us-east-1 \
        --uri 'arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:161262005667:function:events-api/invocations'

#### NOTE
# Retrive previous integration
aws apigateway get-integration --rest-api-id rizogxyxik --resource-id 5w38av --http-method GET

# Create Integration Response for GET /events
aws apigateway put-integration-response --rest-api-id rizogxyxik \
       --resource-id 5w38av --http-method GET \
       --status-code 200 --selection-pattern ""  \
       --region us-east-1

# Create Integration Response for POST /events
aws apigateway put-integration-response --rest-api-id rizogxyxik \
       --resource-id 5w38av --http-method POST \
       --status-code 200 --selection-pattern ""  \
       --region us-east-1

# Add Execution Permision API - Lambda to /events
aws lambda add-permission \
--function-name events-api \
--statement-id 'events_allow_permission' \
--action lambda:InvokeFunction \
--principal apigateway.amazonaws.com \
--source-arn 'arn:aws:execute-api:us-east-1:161262005667:rizogxyxik/*/*/events' \
--region us-east-1


####################################################################
# Deploy de API to stage
# NOTE: Create the backend Lambda
####################################################################

aws apigateway create-deployment --rest-api-id rizogxyxik \
       --region us-east-1 \
       --stage-name test \
       --stage-description 'Test stage' \
       --description 'First deployment'

