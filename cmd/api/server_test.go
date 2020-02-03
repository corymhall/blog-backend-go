package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/corymhall/blog-backend-go/cmd/api/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
	"github.com/rs/zerolog"
)

type mockDynamoDB struct {
	Posts    []*models.Post
	Post     *models.Post
	User     *models.User
	Comments []*models.Comment
	Comment  *models.Comment
	Reply    *models.Reply
	Replies  []*models.Reply
}

func (mdb *mockDynamoDB) DBGetPosts() ([]*models.Post, error) {
	return mdb.Posts, nil
}

func (mdb *mockDynamoDB) DBGetPost(postID string) (*models.Post, error) {
	return mdb.Post, nil
}
func (mdb *mockDynamoDB) DBCreatePost(post *models.Post) error {
	return nil
}

func (mdb *mockDynamoDB) DBUpdatePost(post *models.Post) error {
	return nil
}

func (mdb *mockDynamoDB) DBGetUser(userID string) (*models.User, error) {
	return mdb.User, nil
}

func (mdb *mockDynamoDB) DBCreateUser(user *models.User) error {
	return nil
}

func (mdb *mockDynamoDB) DBGetComments(postID string) ([]*models.Comment, error) {
	return mdb.Comments, nil
}

func (mdb *mockDynamoDB) DBCreateComment(comment *models.Comment) error {
	return nil
}

func (mdb *mockDynamoDB) DBGetReplies(postID, commentID string) ([]*models.Reply, error) {
	return mdb.Replies, nil
}

func (mdb *mockDynamoDB) DBCreateReply(reply *models.Reply) error {
	return nil
}

func newTestHandler(env *Env) http.Handler {
	valv := valve.New()
	baseCtx := valv.Context()
	logger := zerolog.New(os.Stdout)
	logger = logger.With().Timestamp().Logger()

	return chi.ServerBaseContext(
		baseCtx,
		env.Handler(logger),
	)
}

func TestHandlerHealthCheck(t *testing.T) {
	handler := newTestHandler(&Env{&mockDynamoDB{}})
	t.Run("returns health check", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/health", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
		}

		got := response.Body.String()
		want := "healthy"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}

func TestGETPosts(t *testing.T) {
	wantedPost := &models.Post{
		ID:            "1",
		PostText:      "hello",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "location",
		HomeText:      "homeText",
	}
	wantedComment := &models.Comment{
		User:        &models.User{},
		ID:          "1",
		PostID:      "1",
		CommentText: "Hello",
		CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	wanted := []*PostPayload{
		&PostPayload{
			Post: wantedPost,
			Comments: []*CommentPayload{
				&CommentPayload{
					Comment: wantedComment,
				},
			},
		},
	}
	env := &Env{
		db: &mockDynamoDB{
			Posts:    []*models.Post{wantedPost},
			Comments: []*models.Comment{wantedComment},
		},
	}
	handler := newTestHandler(env)
	t.Run("returns posts", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodGet, "/posts", nil)
		res := httptest.NewRecorder()

		var got []*PostPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}

func TestCreatePost(t *testing.T) {
	wantedPost := &models.Post{
		ID:            "1",
		PostText:      "hello",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "location",
		HomeText:      "homeText",
	}
	wanted := &PostPayload{
		Post: wantedPost,
	}
	env := &Env{
		db: &mockDynamoDB{
			Post: wantedPost,
		},
	}

	postPayload := wanted

	jsonPayload, _ := json.Marshal(postPayload)
	handler := newTestHandler(env)
	t.Run("creats a post", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(jsonPayload))
		res := httptest.NewRecorder()

		var got *PostPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusCreated)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}

func TestUpdatePost(t *testing.T) {
	wantedPost := &models.Post{
		ID:            "1",
		PostText:      "hello",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "location",
		HomeText:      "homeText",
	}
	wanted := &PostPayload{
		Post: wantedPost,
	}
	env := &Env{
		db: &mockDynamoDB{
			Post: wantedPost,
		},
	}

	postPayload := wanted

	jsonPayload, _ := json.Marshal(postPayload)
	handler := newTestHandler(env)
	t.Run("updates a post", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(jsonPayload))
		res := httptest.NewRecorder()

		var got *PostPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}

func TestGETPost(t *testing.T) {
	wantedPost := &models.Post{
		ID:            "1",
		PostText:      "hello",
		PostedDate:    time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
		Author:        "author",
		Title:         "title",
		ImageLocation: "location",
		HomeText:      "homeText",
	}
	wantedComment := &models.Comment{
		User:        &models.User{},
		ID:          "1",
		PostID:      "1",
		CommentText: "Hello",
		CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
	}

	wanted := &PostPayload{
		Post: wantedPost,
		Comments: []*CommentPayload{
			&CommentPayload{
				Comment: wantedComment,
			},
		},
	}
	env := &Env{
		db: &mockDynamoDB{
			Post: wantedPost,
			Comments: []*models.Comment{
				wantedComment,
			},
		},
	}
	handler := newTestHandler(env)
	t.Run("returns a single post payload response", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodGet, "/posts/1", nil)
		res := httptest.NewRecorder()

		var got *PostPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}

func TestGETUser(t *testing.T) {
	wantedUser := &models.User{
		ID:          "1",
		DisplayName: "Display Name",
		Email:       "user@gmail.com",
		PhotoURL:    "photourl.com",
		UID:         "12345abcd",
		Role:        "user",
	}

	wanted := &UserPayload{
		User: wantedUser,
	}
	env := &Env{
		db: &mockDynamoDB{
			User: wantedUser,
		},
	}
	handler := newTestHandler(env)
	t.Run("returns user", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodGet, "/user/1", nil)
		res := httptest.NewRecorder()

		var got *UserPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusOK)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})

}

func TestCreateUser(t *testing.T) {
	wantedUser := &models.User{
		ID:          "1",
		DisplayName: "Display Name",
		Email:       "user@gmail.com",
		PhotoURL:    "photourl.com",
		UID:         "12345abcd",
		Role:        "user",
	}

	wanted := &UserPayload{
		User: wantedUser,
	}
	env := &Env{
		db: &mockDynamoDB{
			User: wantedUser,
		},
	}

	userPayload := wanted

	jsonPayload, _ := json.Marshal(userPayload)
	handler := newTestHandler(env)
	t.Run("creats a user", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonPayload))
		res := httptest.NewRecorder()

		var got *UserPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusCreated)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})

}

func TestCreateComment(t *testing.T) {
	wantedComment := &models.Comment{
		User:        &models.User{},
		ID:          "1",
		PostID:      "1",
		CommentText: "Hello",
		CommentDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	wanted := &CommentPayload{
		Comment: wantedComment,
	}
	env := &Env{
		db: &mockDynamoDB{
			Comment: wantedComment,
		},
	}

	commentPayload := wanted

	jsonPayload, _ := json.Marshal(commentPayload)
	handler := newTestHandler(env)
	t.Run("creats a comment", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(jsonPayload))
		res := httptest.NewRecorder()

		var got *CommentPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusCreated)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}

func TestCreateReply(t *testing.T) {
	wantedReply := &models.Reply{
		User:      &models.User{},
		ID:        "1:1",
		ReplyText: "Hello",
		ReplyDate: time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	wanted := &ReplyPayload{
		Reply: wantedReply,
	}
	env := &Env{
		db: &mockDynamoDB{
			Reply: wantedReply,
		},
	}

	replyPayload := wanted

	jsonPayload, _ := json.Marshal(replyPayload)
	handler := newTestHandler(env)
	t.Run("creats a comment", func(t *testing.T) {
		rq, _ := http.NewRequest(http.MethodPost, "/replies", bytes.NewBuffer(jsonPayload))
		res := httptest.NewRecorder()

		var got *ReplyPayload

		handler.ServeHTTP(res, rq)

		if status := res.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got '%v' want '%v'", status, http.StatusCreated)
		}

		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Post, '%v'", res.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got '%s', want '%s'", got, wanted)
		}

	})
}
