package models

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestDBGetComments(t *testing.T) {
	items := []map[string]*dynamodb.AttributeValue{}
	attribute := Comment{

		ID:          "1",
		CommentText: "hello world",
		CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		PostID:      "1",
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	items = append(items, item)
	cases := []struct {
		Resp     dynamodb.QueryOutput
		Expected []*Comment
	}{
		{
			Resp: dynamodb.QueryOutput{
				Items: items,
			},
			Expected: []*Comment{
				&Comment{
					ID:          "1",
					CommentText: "hello world",
					CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
					PostID:      "1",
					//User:          nil},

				},
			},
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedQuery{Resp: c.Resp},
		}
		items, err := d.DBGetComments("1")
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		for i, p := range items {
			if !compareComment(p, c.Expected[i]) {
				t.Errorf("expected %v message, got %v", p, c.Expected[i])
			}
		}
	}
}

func TestDBCreateComment(t *testing.T) {
	attribute := Comment{

		ID:          "1",
		CommentText: "hello world",
		CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		PostID:      "1",
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	cases := []struct {
		Resp     dynamodb.PutItemOutput
		Expected *Comment
	}{
		{
			Resp: dynamodb.PutItemOutput{
				Attributes: item,
			},
			Expected: nil,
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedPutItem{Resp: c.Resp},
		}
		err := d.DBCreateComment(&attribute)
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
	}
}

func compareComment(a, b *Comment) bool {
	if a.ID != b.ID {
		return false
	}
	if a.CommentText != b.CommentText {
		return false
	}
	if a.CommentDate != b.CommentDate {
		return false
	}
	return true
}
