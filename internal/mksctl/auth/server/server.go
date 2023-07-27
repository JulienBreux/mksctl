package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/coreos/go-oidc"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const (
	handlerPath           = "/sso-callback"
	handlerCodeQueryParam = "code"

	readHeaderTimeout = 3 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// Server represents a server interface
type Server interface {
	Run() error
}

// New creates a new instance of server
func New() (Server, error) {
	provider, err := oidc.NewProvider(context.Background(), issuerURL())
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:    config.Config.AuthClientID,
		RedirectURL: redirectURL(),
		Endpoint:    provider.Endpoint(),
	}

	srv := &http.Server{
		Addr:              serverAddress(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &server{
		interrupt:    make(chan any, 1),
		oauth2Config: oauth2Config,
		server:       srv,
	}, nil
}

type server struct {
	interrupt    chan any
	oauth2Config oauth2.Config
	server       *http.Server
}

// Run starts the server
func (s *server) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	if err := openBrowser(authURL()); err != nil {
		return err
	}

	return s.stop()
}

// start the callback server
func (s *server) start() error {
	http.HandleFunc(handlerPath, s.handlerFunc)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

// stop the callback server
func (s *server) stop() error {
	stopFromOS(s.interrupt)

	<-s.interrupt

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// handlerFunc returns the callback handler
func (s *server) handlerFunc(w http.ResponseWriter, r *http.Request) {
	// Check code query parametter
	code := r.URL.Query().Get(handlerCodeQueryParam)
	if code == "" {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "Code query parametter is missing.")
		return
	}

	// Converts an authorization code into a token.
	oauth2Token, err := s.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Unable to exchange the authorization code for a token.\n\nError: %v", err)
		return
	}

	// Read access token
	rawAccessToken, ok := oauth2Token.Extra("access_token").(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Access token not found.")
		return
	}

	// Save access token to configuration
	// TODO: Save to keyring
	config.Config.AuthAccessToken = rawAccessToken
	if err := config.Save(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Unable to save configuration file.")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "You can close this page. You're now authenticated.")

	s.interrupt <- true
}

func authURL() string {
	return fmt.Sprintf(
		"%v/realms/%v/protocol/openid-connect/auth?client_id=%v&redirect_uri=%v&response_type=code",
		config.Config.AuthServerURL,
		config.Config.AuthClientRealm,
		config.Config.AuthClientID,
		url.QueryEscape(redirectURL()),
	)
}

func issuerURL() string {
	return fmt.Sprintf("%s/realms/%s", config.Config.AuthServerURL, config.Config.AuthClientRealm)
}

func serverAddress() string {
	return fmt.Sprintf("localhost:%v", config.Config.AuthCallbackPort)
}

func redirectURL() string {
	return fmt.Sprintf("http://%s/%s", serverAddress(), handlerPath)
}

func openBrowser(authURL string) error {
	if err := browser.OpenURL(authURL); err != nil {
		return err
	}
	fmt.Printf("Your browser has been opened to visit: \n\n\t%s\n\n", authURL)
	return nil
}

func stopFromOS(stop chan any) {
	stopOS := make(chan os.Signal, 1)
	go func() {
		<-stopOS
		stop <- true
	}()
	signal.Notify(stopOS, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}
