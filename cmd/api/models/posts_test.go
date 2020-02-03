package models

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestDBGetPost(t *testing.T) {
	items := []map[string]*dynamodb.AttributeValue{}
	attribute := Post{
		ID:            "1",
		PostText:      "hello world",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "abcd",
		HomeText:      "home text",
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	items = append(items, item)
	cases := []struct {
		Resp     dynamodb.QueryOutput
		Expected *Post
	}{
		{
			Resp: dynamodb.QueryOutput{
				Items: items,
			},
			Expected: &Post{
				ID:            "1",
				PostText:      "hello world",
				PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
				Author:        "author",
				Title:         "title",
				ImageLocation: "abcd",
				HomeText:      "home text",
				//User:          nil},

			},
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedQuery{Resp: c.Resp},
		}
		items, err := d.DBGetPost("1")
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if !comparePost(items, c.Expected) {
			t.Errorf("expected %v message, got %v", items, c.Expected)
		}
	}
}

func TestDBGetPosts(t *testing.T) {
	items := []map[string]*dynamodb.AttributeValue{}
	attribute := Post{

		ID:            "1",
		PostText:      "hello world",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "abcd",
		HomeText:      "home text",
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	items = append(items, item)
	cases := []struct {
		Resp     dynamodb.QueryOutput
		Expected []*Post
	}{
		{
			Resp: dynamodb.QueryOutput{
				Items: items,
			},
			Expected: []*Post{
				&Post{
					ID:            "1",
					PostText:      "hello world",
					PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
					Author:        "author",
					Title:         "title",
					ImageLocation: "abcd",
					HomeText:      "home text",
					//User:          nil},

				},
			},
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedQuery{Resp: c.Resp},
		}
		items, err := d.DBGetPosts()
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		for i, p := range items {
			if !comparePost(p, c.Expected[i]) {
				t.Errorf("expected %v message, got %v", p, c.Expected[i])
			}
		}
	}
}

func TestDBCreatePost(t *testing.T) {
	attribute := Post{
		ID:            "1",
		PostText:      "hello world",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "abcd",
		HomeText:      "home text",
		//User:          nil}
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	cases := []struct {
		Resp     dynamodb.PutItemOutput
		Expected *Post
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
		err := d.DBCreatePost(&attribute)
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
	}
}

func comparePost(a, b *Post) bool {
	if a.Author != b.Author {
		return false
	}
	if a.HomeText != b.HomeText {
		return false
	}
	if a.ID != b.ID {
		return false
	}
	if a.ImageLocation != b.ImageLocation {
		return false
	}
	if a.PostText != b.PostText {
		return false
	}
	if a.Title != b.Title {
		return false
	}
	return true
}
