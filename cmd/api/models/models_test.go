package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedQuery struct {
	dynamodbiface.DynamoDBAPI
	Resp dynamodb.QueryOutput
}

type mockedScan struct {
	dynamodbiface.DynamoDBAPI
	Resp dynamodb.ScanOutput
}

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Resp dynamodb.PutItemOutput
}

type mockedGetItem struct {
	dynamodbiface.DynamoDBAPI
	Resp dynamodb.GetItemOutput
}

func (m mockedQuery) Query(in *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return &m.Resp, nil
}

func (m mockedScan) Scan(in *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &m.Resp, nil
}

func (m mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &m.Resp, nil
}

func (m mockedGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &m.Resp, nil
}
