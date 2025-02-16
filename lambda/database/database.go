package database

import (
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/session"
)

type DynamoDBClient struct {
  databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
  dbSession := session.Must(session.NewSession())
  db := dynamodb.New(dbSession)
  return DynamoDBClient{
    databaseStore: db,
  }
}
