package cmd

import (
	"github.com/coreos/go-oidc"
	"net/http"
)

type app struct {
	clientID       string
	clientSecret   string
	redirectURI    string
	state          string
	tokenRetrieved chan int

	verifier *oidc.IDTokenVerifier
	provider *oidc.Provider

	// Does the provider use "offline_access" scope to request a refresh token
	// or does it use "access_type=offline" (e.g. Google)?
	offlineAsScope bool

	client *http.Client

	kubeconfig *KubeConfig
}

//Clusters struct represents type of Cluster.
type Clusters struct {
	Name         string
	Address      string
	CAServerName string
	Certificate  string
	Tier         string
}

//Tiers struct represents type of Tier.
type Tiers struct {
	Name         string
	Issuer       string
	CAServerName string
}
