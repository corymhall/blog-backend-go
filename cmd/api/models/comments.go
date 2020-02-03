package models

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Comment struct {
	User        *User     `json:"user"`
	ID          string    `json:"id"`
	PostID      string    `json:"post_id"`
	CommentText string    `json:"comment_text"`
	CommentDate time.Time `json:"comment_date"`
}

func (db *DB) DBGetComments(postID string) ([]*Comment, error) {
	comments := []*Comment{}
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(postID),
			},
		},
		KeyConditionExpression: aws.String("post_id = :v1"),
		//ProjectionExpression:   aws.String("comment_date"),
		TableName: aws.String("Comments"),
	}

	result, err := db.Svc.Query(input)
	if err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *DB) DBCreateComment(comment *Comment) error {
	item, err := dynamodbattribute.MarshalMap(comment)
	if err != nil {
		return err
	}
	_, err = db.Svc.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Comments"),
	})
	if err != nil {
		return err
	}
	return nil
}
