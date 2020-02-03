package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Reply struct {
	User      *User     `json:"user"`
	ID        string    `json:"id"`
	ReplyText string    `json:"reply_text"`
	ReplyDate time.Time `json:"reply_date"`
}

func (db *DB) DBGetReplies(postID, commentID string) ([]*Reply, error) {
	replies := []*Reply{}
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(fmt.Sprintf("%s#%s", postID, commentID)),
			},
		},
		KeyConditionExpression: aws.String("id = :v1"),
		TableName:              aws.String("Reply"),
	}

	result, err := db.Svc.Query(input)
	if err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &replies); err != nil {
		return nil, err
	}

	return replies, nil
}

func (db *DB) DBCreateReply(reply *Reply) error {
	item, err := dynamodbattribute.MarshalMap(reply)
	if err != nil {
		return err
	}
	_, err = db.Svc.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Reply"),
	})
	if err != nil {
		return err
	}
	return nil
}
