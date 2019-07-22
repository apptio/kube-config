package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	t "github.com/apptio/kube-config/pkg/templates"
	"github.com/coreos/go-oidc"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/oauth2"
)

func startHttpServer(a *app, listen string) *http.Server {

	u, err := url.Parse(a.redirectURI)
	if err != nil {
		log.Fatalf("parse redirect-uri: %v", err)
	}
	log.Debug("redirect-uri: ", a.redirectURI)

	listenURL, err := url.Parse(listen)
	if err != nil {
		log.Fatalf("parse listen address: %v", err)
	}

	log.Debug("listenURL: ", listenURL)

	srv := &http.Server{Addr: listenURL.Host}

	log.Debug("Starting http server on: ", listenURL.Host)

	http.HandleFunc("/", a.handleLogin)
	http.HandleFunc(u.Path, a.handleCallback)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Infof("Closing HTTP Server: %s", err)
		}
	}()

	return srv

}
func (a *app) handleCallback(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token *oauth2.Token
	)

	ctx := oidc.ClientContext(r.Context(), a.client)
	oauth2Config := a.oauth2Config(nil)
	switch r.Method {
	case "GET":
		// Authorization redirect callback from OAuth2 auth flow.
		if errMsg := r.FormValue("error"); errMsg != "" {
			http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
			return
		}
		if state := r.FormValue("state"); state != a.state {
			http.Error(w, fmt.Sprintf("expected state %q got %q", a.state, state), http.StatusBadRequest)
			return
		}
		token, err = oauth2Config.Exchange(ctx, code)
	case "POST":
		// Form request from frontend to refresh a token.
		refresh := r.FormValue("refresh_token")
		if refresh == "" {
			http.Error(w, fmt.Sprintf("no refresh_token in request: %q", r.Form), http.StatusBadRequest)
			return
		}
		t := &oauth2.Token{
			RefreshToken: refresh,
			Expiry:       time.Now().Add(-time.Hour),
		}
		token, err = oauth2Config.TokenSource(ctx, t).Token()
	default:
		http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		return
	}

	idToken, err := a.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify ID token: %v", err), http.StatusInternalServerError)
		return
	}
	var claims json.RawMessage
	idToken.Claims(&claims)

	buff := new(bytes.Buffer)
	json.Indent(buff, []byte(claims), "", "  ")

	t.RenderToken(w, a.redirectURI, rawIDToken, token.RefreshToken, buff.Bytes())

	a.kubeconfig.Generate(rawIDToken, token.RefreshToken)
	a.tokenRetrieved <- 1

}

func (a *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	scopes := []string{"groups"}
	clients := []string{a.clientID}
	for _, client := range clients {
		scopes = append(scopes, "audience:server:client_id:"+client)
	}

	authCodeURL := ""

	// FIXME: add ability to specify offline access on commandline
	scopes = append(scopes, "openid", "profile", "email", "offline_access")

	authCodeURL = a.oauth2Config(scopes).AuthCodeURL(a.state, oauth2.AccessTypeOffline)

	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}
