aws iam create-role \
--role-name 'PersonsAPILambdaRole' \
--assume-role-policy-document 'file://assume-role-policy-document.json' \
--description 'Allows Lambda functions to call AWS services on your behalf.'
# {
#     "Role": {
#         "AssumeRolePolicyDocument": {
#             "Version": "2012-10-17",
#             "Statement": [
#                 {
#                     "Action": "sts:AssumeRole",
#                     "Effect": "Allow",
#                     "Principal": {
#                         "Service": "lambda.amazonaws.com"
#                     }
#                 }
#             ]
#         },
#         "RoleId": "AROAJBA47PFUG4XP4CDNK",
#         "CreateDate": "2018-08-22T23:12:29Z",
#         "RoleName": "PersonsAPILambdaRole",
#         "Path": "/",
#         "Arn": "arn:aws:iam::161262005667:role/PersonsAPILambdaRole"
#     }
# }

# Attach managed policy AWSLambdaBasicExecutionRole to the Rol
aws iam attach-role-policy \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
--role-name PersonsAPILambdaRole

# Create the dynamodb table
./create_table.sh

# Embed inline policy to Role so it can update the databse
aws iam put-role-policy \
--role-name PersonsAPILambdaRole \
--policy-name DynamoDBPolicy \
--policy-document file://dynamodb-inline-policy.json