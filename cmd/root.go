// Copyright © 2019 Apptio
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code used from the dex example app is licensed under the original Apache 2.0 license
// https://github.com/dexidp/dex/blob/master/cmd/example-app

package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"context"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/jaxxstorm/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/skratchdot/open-golang/open"
	//"github.com/davecgh/go-spew/spew"
)

var (
	cfgFile        string
	a              app
	issuerURL      string
	listen         string
	debug          bool
	userName       string
	namespace      string
	clusters       []Clusters
	cluster        string
	outputFilePath string
	insecure       bool
	tiers          []Tiers
	tier           string
	listTiers      bool
)

// for the version command
var Version string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kube-config",
	Short: "Configure a kubeconfig for Kubernetes Clusters that have Dex as an IDP",
	Long: `kube-config is a method of authenticating to a Kubernetes cluster which uses Dex an IDP.
	
The tool will authenticate you using your credentials provider, and grab a token.
It will also generate you a Kubernetes configuration file based on your login credentials.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// turn on debug logging
		if debug {

			//log.SetLevel(log.DebugLevel)

			if a.client == nil {
				a.client = &http.Client{
					Transport: debugTransport{http.DefaultTransport},
				}
			} else {
				a.client.Transport = debugTransport{a.client.Transport}
			}
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: insecure,
		}

		// set the default http client
		if a.client == nil {
			c := cleanhttp.DefaultClient()
			t := cleanhttp.DefaultTransport()
			t.TLSClientConfig = tlsConfig
			c.Transport = t
			a.client = c
		}

		err := viper.UnmarshalKey("tiers", &tiers)

		if listTiers {

			var entry []string

			log.Info("Outputting available tiers")

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Issuer"})

			for _, v := range tiers {
				entry = append(entry, v.Name)
				entry = append(entry, v.Issuer)
				table.Append(entry)
				entry = entry[:0]
			}

			table.Render()

			os.Exit(0)

		}

		for _, i := range tiers {
			if i.Name == tier {
				issuerURL = i.Issuer
				log.Debug("Using Issuer URL: ", issuerURL)
			}
		}

		ctx := oidc.ClientContext(context.Background(), a.client)
		provider, err := oidc.NewProvider(ctx, issuerURL)
		if err != nil {
			log.Fatalf("Failed to query provider %q: %v", issuerURL, err)
		}

		var s struct {
			// What scopes does a provider support?
			//
			// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
			ScopesSupported []string `json:"scopes_supported"`
		}

		if err := provider.Claims(&s); err != nil {
			log.Fatalf("Failed to parse provider scopes_supported: %v", err)
		}

		if len(s.ScopesSupported) == 0 {
			// scopes_supported is a "RECOMMENDED" discovery claim, not a required
			// one. If missing, assume that the provider follows the spec and has
			// an "offline_access" scope.
			a.offlineAsScope = true
		} else {
			// See if scopes_supported has the "offline_access" scope.
			a.offlineAsScope = func() bool {
				for _, scope := range s.ScopesSupported {
					if scope == oidc.ScopeOfflineAccess {
						return true
					}
				}
				return false
			}()

		}

		a.provider = provider
		a.verifier = provider.Verifier(&oidc.Config{ClientID: a.clientID})
		a.state, err = generateState()
		a.tokenRetrieved = make(chan int, 1)

		if err != nil {
			log.Fatal("Error generating state: ", err)
		}

		fileHandle := os.Stdout

		if outputFilePath != "" {
			var ok bool

			log.Debug(outputFilePath)
			_, err := os.Stat(outputFilePath)

			if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
				ok = prompt.Confirm("The kube-config directory doesn't exist at %s - Would you like to create it?: (y/n) ", outputFilePath)
				if ok {
					err := os.MkdirAll(outputFilePath, 0700)
					if err != nil {
						log.Fatalf("Error creating directory: %s", err)
					}
				} else {
					log.Fatalf("Cannot continue, please create directory %s", outputFilePath)
				}
			}

			fileHandle, err = os.Create(outputFilePath + "/" + tier + "-config.yml")
			if err != nil {
				log.Fatalf("Error creating kubeconfig: %v", err)
			}
		}

		userName = viper.GetString("username")
		a.clientSecret = viper.GetString("client-secret")

		err = viper.UnmarshalKey("clusters", &clusters)

		if err != nil {
			log.Fatalf("Error reading datacenters: %v", err)
		}

		if clusters == nil {
			log.Fatal("No clusters specified in config file")
		}

		if userName == "" {
			log.Fatal("Please specify a username to login with")
		}

		if a.clientSecret == "" {
			log.Fatal("Client secret must be specified in config file")
		}

		// loop through the clusters
		// and grab the clusters in the tier we specified
		tierClusters := clusters[:0]

		for _, c := range clusters {
			if c.Tier == tier {
				tierClusters = append(tierClusters, c)
			}
		}

		kubeConfig, err := NewKubeConfig(cluster, tierClusters, userName, namespace, fileHandle, a.clientID, issuerURL, a.clientSecret)

		if err != nil {
			log.Warn("Error generating KubeConfig: ", err)
		}

		a.kubeconfig = kubeConfig

		srv := startHttpServer(&a, listen)

		open.Run(listen)

		// block
		<-a.tokenRetrieved

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Info("Error shutting down the server gracefully - this isn't usually an issue: ", err) // failure/timeout shutting down the server gracefully
		}

		log.Infof("Your kubeconfig has been written to %s/%s-config.yml", outputFilePath, tier)
		log.Warnf("You must set the KUBECONFIG environment variable. It's suggested you export all the files in %s. You can do this by running:", outputFilePath)
		log.Warnf("export KUBECONFIG=$HOME/.kube/config; for i in %s/*.yml; do export KUBECONFIG=$KUBECONFIG:$i; done", outputFilePath)

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	log.SetLevel(log.InfoLevel)

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&a.clientID, "client-id", "kube-config", "OAuth2 client ID of this application")
	RootCmd.PersistentFlags().StringVar(&a.redirectURI, "redirect-uri", "http://127.0.0.1:5555/callback", "Callback URL for OAuth2 responses")
	RootCmd.PersistentFlags().StringVarP(&tier, "tier", "t", "dev", "Tier to authenticate for")
	RootCmd.PersistentFlags().StringVar(&issuerURL, "issuer", "http://localhost", "URL of the OpenID Connect issuer")
	RootCmd.PersistentFlags().StringVar(&listen, "listen", "http://127.0.0.1:5555", "HTTP(S) address to listen at")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Print all request and responses from the OpenID Connect issuer")
	RootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Username for login")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Set your default namespace across ALL contexts in the tier")
	RootCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "", "Cluster to use as current-context (default is first in config)")
	RootCmd.PersistentFlags().StringVarP(&outputFilePath, "output", "o", os.Getenv("HOME")+"/.kube/config.d", "Path to write Kubeconfig file")
	RootCmd.PersistentFlags().BoolVarP(&insecure, "no-verify-ssl", "k", false, "If specified, disable SSL cert checking (WARNING: UNSAFE)")
	RootCmd.PersistentFlags().BoolVar(&listTiers, "list-tiers", false, "If specified, the program will list the available tiers and then exit")
	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	RootCmd.PersistentFlags().MarkHidden("client-id")
	RootCmd.PersistentFlags().MarkHidden("listen")
	RootCmd.PersistentFlags().MarkHidden("redirect-uri")
	RootCmd.PersistentFlags().MarkHidden("issuer")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	viper.SetConfigName("kube-config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kube-config") // adding home directory as first search path
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/etc")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file: ", viper.ConfigFileUsed())
	}

	currentUser, err := user.Current()

	if err != nil {
		log.Debug("Error retrieving current user: ", err)
	}

	viper.SetDefault("username", currentUser.Username)
}
