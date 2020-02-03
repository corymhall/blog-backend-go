package models

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestDBGetRepies(t *testing.T) {
	items := []map[string]*dynamodb.AttributeValue{}
	attribute := Reply{

		ID:        "1#1",
		ReplyText: "hello world",
		ReplyDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	items = append(items, item)
	cases := []struct {
		Resp     dynamodb.QueryOutput
		Expected []*Reply
	}{
		{
			Resp: dynamodb.QueryOutput{
				Items: items,
			},
			Expected: []*Reply{
				&Reply{
					ID:        "1#1",
					ReplyText: "hello world",
					ReplyDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
					//User:          nil},

				},
			},
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedQuery{Resp: c.Resp},
		}
		items, err := d.DBGetReplies("1", "1")
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		for i, p := range items {
			if !compareReply(p, c.Expected[i]) {
				t.Errorf("expected %v message, got %v", p, c.Expected[i])
			}
		}
	}
}

func TestDBCreateReply(t *testing.T) {
	attribute := Reply{

		ID:        "1#1",
		ReplyText: "hello world",
		ReplyDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	cases := []struct {
		Resp     dynamodb.PutItemOutput
		Expected *Reply
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
		err := d.DBCreateReply(&attribute)
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
	}
}

func compareReply(a, b *Reply) bool {
	if a.ID != b.ID {
		return false
	}
	if a.ReplyText != b.ReplyText {
		return false
	}
	if a.ReplyDate != b.ReplyDate {
		return false
	}
	return true
}
