package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/corymhall/blog-backend-go/cmd/api/models"
	ht "github.com/corymhall/blog-backend-go/pkg/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/valve"
	"github.com/rs/zerolog"
)

type Env struct {
	db models.Datastore
}

func main() {
	var (
		httpAddr = flag.String("http.addr", fmt.Sprintf(":%s", os.Getenv("PORT")), "HTTP listen address")
	)

	flag.Parse()
	valv := valve.New()
	baseCtx := valv.Context()
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	env := &Env{db: &models.DB{
		Svc: dynamodb.New(sess),
	}}

	var logger zerolog.Logger
	{
		logger = zerolog.New(os.Stdout)
		logger = logger.With().
			Timestamp().
			Logger()
	}

	var srv http.Server

	logger = logger.With().Str("transport", "http").Logger()
	logger.Info().Str("addr", *httpAddr).Msg("starting http server")

	// server config
	{
		srv.Addr = *httpAddr
		srv.Handler = chi.ServerBaseContext(
			baseCtx,
			env.Handler(logger),
		)
		srv.ReadTimeout = time.Second * 30
		srv.WriteTimeout = time.Second * 30
	}
	logger.Info().Msg("preparing to listen to requests")
	srv.ListenAndServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Info().Msg("shutting down...")

			// first valve
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			srv.Shutdown(ctx)

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				logger.Info().Msg("not all connections done")
			case <-ctx.Done():

			}
		}
	}()
}

func (env *Env) Handler(logger zerolog.Logger) http.Handler {

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://www.pleasantplacesblog.com", "https://pleasantplacesblog.com"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(ht.NewLogger(logger))
	r.Use(middleware.RequestID)
	/*r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.RequestIDHandler("req_id","Request-Id"))
	r.Use(hlog.MethodHandler("method"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.URLHandler("url"))*/

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	r.Mount("/posts", postRouter(env))
	r.Mount("/user", userRouter(env))
	r.Mount("/comments", commentRouter(env))
	r.Mount("/replies", replyRouter(env))
	return r

}
