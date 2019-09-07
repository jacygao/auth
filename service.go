// Package auth contains a service that handles authentication.
package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Config defines configuration variables.
type Config struct {
	HTTPPort                      int
	VerificationCodeExpiryMinutes int
}

// Service defines Auth Service instance structure and dependecies.
type Service struct {
	conf     Config
	mailer   Mailer
	handler  *Handler
	callback Callback
}

// Start starts the Auth Service.
func (s *Service) Start() {
	router := httprouter.New()
	router.GET("/health", s.handler.Health)
	router.GET("/user/find/:id", s.handler.User)
	router.PUT("/user/forget/:id", s.handler.Forget)
	router.PUT("/user/forget", s.handler.SetPassword)
	router.POST("/user/login", s.handler.Login)
	router.POST("/user/signup", s.handler.Signup)
	router.POST("/user/verify", s.handler.Verify)

	router.GET("/authenticate/:token", s.handler.Authenticate)
	router.GET("/profile/:id", s.handler.Profile)

	// TODO: Consolidate this
	handler := cors.AllowAll().Handler(router)

	log.Printf("Authentication Server listenning on port %d", s.conf.HTTPPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.conf.HTTPPort), handler))
}

// Stop contains processes for a graceful shutdown.
func (s *Service) Stop() {
	// Close everything
	log.Println("Shutting down...")
}

// WithMongoDB registers Mongo DB in the Auth Service instance.
func (s *Service) WithMongoDB(session *mgo.Session, dbName, collection string) (*Service, error) {
	model, err := NewMongoDB(session, dbName, collection)
	if err != nil {
		return nil, err
	}
	if s.mailer == nil {
		return nil, errors.New("you must call withMailer first")
	}
	s.handler = NewHandler(s.conf, model, s.mailer, s.callback)
	return s, nil
}

func (s *Service) WithMongoDriver(client *mongo.Client, dbName, collection string) (*Service, error) {
	model, err := NewMongo(client, dbName, collection)
	if err != nil {
		return nil, err
	}
	if s.mailer == nil {
		return nil, errors.New("you must call withMailer first")
	}
	s.handler = NewHandler(s.conf, model, s.mailer, s.callback)
	return s, nil
}

// WithMailer registers a mail service in the Auth Service instance.
// TODO: Remove Mailer as it should be handled in the Callback methods.
func (s *Service) WithMailer(ms Mailer) *Service {
	s.mailer = ms
	return s
}

func (s *Service) WithCallback(cb Callback) *Service {
	s.callback = cb
	return s
}

// NewService initialises a new Auth Service instance.
func NewService(c Config) *Service {
	// Default verification expiry to 10 minutes
	if c.VerificationCodeExpiryMinutes == 0 {
		c.VerificationCodeExpiryMinutes = 10
	}
	return &Service{
		conf: c,
	}
}
