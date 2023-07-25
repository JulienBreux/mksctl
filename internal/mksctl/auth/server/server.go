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
	ssoCallbackPath   = "/sso-callback"
	ssoCodeQueryParam = "code"
)

// Server represents a server interface
type Server interface {
	Run() error
}

// New creates a new instance of server
func New() Server {
	return &server{}
}

type server struct{}

// Run starts the server
func (s *server) Run() error {
	stop := make(chan any, 1)

	// Configure callback webserver
	const readHeaderTimeout = 3 * time.Second
	server := &http.Server{
		Addr:              serverAddress(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Callback handler
	ssoCallbackFunc, err := ssoCallback(stop)
	if err != nil {
		return err
	}
	http.HandleFunc(ssoCallbackPath, ssoCallbackFunc)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Open browers
	if err := browser.OpenURL(authURL()); err != nil {
		return err // TODO: Clean error message
	}

	fmt.Printf("Your browser has been opened to visit: \n\n\t%s\n\n", authURL())

	// Catch system signal
	stopFromOS(stop)

	// Stop in server gracefully
	<-stop
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err // TODO: Clean error message
	}

	// fmt.Printf("\nCommand killed by keyboard interrupt.\n")

	return nil
}

func ssoCallback(stop chan<- any) (http.HandlerFunc, error) {
	// Configure OIDC provider
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerURL())
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:    config.Config.AuthClientID,
		RedirectURL: redirectURL(),
		Endpoint:    provider.Endpoint(),
	}

	callback := func(w http.ResponseWriter, r *http.Request) {
		// Check code query parametter
		code := r.URL.Query().Get(ssoCodeQueryParam)
		if code == "" {
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprintf(w, "Code query parametter is missing.")
			return
		}

		// Converts an authorization code into a token.
		oauth2Token, err := oauth2Config.Exchange(ctx, code)
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

		stop <- true
	}

	return callback, nil
}

func stopFromOS(stop chan any) {
	stopOS := make(chan os.Signal, 1)
	go func() {
		<-stopOS
		stop <- true
	}()
	signal.Notify(stopOS, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
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
	return fmt.Sprintf("http://%s/%s", serverAddress(), ssoCallbackPath)
}
