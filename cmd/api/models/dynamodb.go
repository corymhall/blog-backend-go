package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DB struct {
	Svc dynamodbiface.DynamoDBAPI
}

type Datastore interface {
	DBGetPosts() ([]*Post, error)
	DBGetPost(postID string) (*Post, error)
	DBCreatePost(post *Post) error
	DBUpdatePost(post *Post) error
	DBGetUser(userID string) (*User, error)
	DBCreateUser(user *User) error
	DBGetComments(postID string) ([]*Comment, error)
	DBGetReplies(postID, commentID string) ([]*Reply, error)
	DBCreateComment(comment *Comment) error
	DBCreateReply(reply *Reply) error
}
