package models

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestDBGetUser(t *testing.T) {
	attribute := User{
		ID:          "1",
		DisplayName: "User 1",
		Email:       "user@email.com",
		PhotoURL:    "/images/somepath.jpg",
		ProviderID:  "google.com",
		Role:        "Reader",
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	cases := []struct {
		Resp     dynamodb.GetItemOutput
		Expected *User
	}{
		{
			Resp: dynamodb.GetItemOutput{
				Item: item,
			},
			Expected: &User{
				ID:          "1",
				DisplayName: "User 1",
				Email:       "user@email.com",
				PhotoURL:    "/images/somepath.jpg",
				ProviderID:  "google.com",
				Role:        "Reader",
			},
		},
	}

	for _, c := range cases {
		d := DB{
			Svc: mockedGetItem{Resp: c.Resp},
		}
		items, err := d.DBGetUser("1")
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if !compareUser(items, c.Expected) {
			t.Errorf("expected %v message, got %v", items, c.Expected)
		}
	}
}

func TestDBCreateUser(t *testing.T) {
	attribute := User{
		ID:          "1",
		DisplayName: "User 1",
		Email:       "user@email.com",
		PhotoURL:    "/images/somepath.jpg",
		ProviderID:  "google.com",
		Role:        "Reader",
	}
	item, _ := dynamodbattribute.MarshalMap(attribute)
	cases := []struct {
		Resp     dynamodb.PutItemOutput
		Expected *User
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
		err := d.DBCreateUser(&attribute)
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
	}
}

func compareUser(a, b *User) bool {
	if a.ID != b.ID {
		return false
	}
	if a.DisplayName != b.DisplayName {
		return false
	}
	if a.Email != b.Email {
		return false
	}
	if a.ProviderID != b.ProviderID {
		return false
	}
	if a.PhotoURL != b.PhotoURL {
		return false
	}
	if a.Role != b.Role {
		return false
	}
	return true
}
