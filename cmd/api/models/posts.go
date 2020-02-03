package models

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Post struct {
	ID            string    `json:"id"`
	PostText      string    `json:"post_text"`
	PostedDate    time.Time `json:"posted_date"`
	Author        string    `json:"author"`
	Title         string    `json:"title"`
	ImageLocation string    `json:"image_location"`
	HomeText      string    `json:"home_text"`
	User          *User     `json:"user,omitempty"`
}

func (db *DB) DBGetPosts() ([]*Post, error) {
	posts := []*Post{}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String("Pleasant Places"),
			},
		},
		KeyConditionExpression: aws.String("author = :v1"),
		IndexName:              aws.String("author-posted_date-index"),
		TableName:              aws.String("Posts"),
		ScanIndexForward:       aws.Bool(false),
	}
	res, err := db.Svc.Query(input)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &posts); err != nil {
		return nil, err
	}

	return posts, nil

}

func (db *DB) DBGetPost(postID string) (*Post, error) {
	var post Post

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(postID),
			},
		},
		KeyConditionExpression: aws.String("id = :v1"),
		TableName:              aws.String("Posts"),
	}

	res, err := db.Svc.Query(input)
	if err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalMap(res.Items[0], &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (db *DB) DBCreatePost(post *Post) error {
	item, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return err
	}
	_, err = db.Svc.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Posts"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DBUpdatePost(post *Post) error {
	return nil
}
