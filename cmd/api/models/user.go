package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	PhotoURL    string `json:"photo_url"`
	ProviderID  string `json:"provider_id"`
	UID         string `json:"uid"`
	Role        string `json:"role"`
}

func (db *DB) DBGetUser(userID string) (*User, error) {
	var user User

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(userID),
			},
		},
		TableName: aws.String("User"),
	}

	res, err := db.Svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalMap(res.Item, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) DBCreateUser(user *User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}
	_, err = db.Svc.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("User"),
	})
	if err != nil {
		return err
	}
	return nil
}
