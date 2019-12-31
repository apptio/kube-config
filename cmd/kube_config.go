package cmd

import (
	"io"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

// KubeConfig holds the information necessary to generate a Kubernetes configuration file which icludes the server's CA, the api url and where to write the file to.
type KubeConfig struct {
	Cluster      string
	Clusters     []Clusters
	Username     string
	NS           string
	CAServerName string
	tmpl         *template.Template
	Output       io.ReadWriteCloser
	ClientID     string
	Issuer       string
	ClientSecret string
	Tier         string
}

var content = `
apiVersion: v1
clusters:
{{- $caservername := .CAServerName }}
{{- range .Clusters}}
- cluster:
    {{- if .Certificate }}
    certificate-authority-data: {{.Certificate}}
    {{- else }}{{- $sname := or .CAServerName $caservername }}
    certificate-authority-data: {{printf "https://%s:6444/ca.crt" .Address | getCert $sname }}
    {{- end }}
    server: https://{{.Address}}:6443
  name: {{.Name}}
{{- end}}
contexts:
{{- range .Clusters}}
- context:
    cluster: {{.Name}}
    user: {{$.Username}}@{{.Tier}}
    {{- if $.NS }}
    namespace: {{$.NS}}
    {{- end }}
  name: {{$.Username}}@{{.Name}}
{{- end}}
current-context: {{.Username}}@
	{{- if $.Cluster }}
		{{- $.Cluster }}
	{{- else }}
		{{- (index .Clusters 0).Name }}
	{{- end }}
kind: Config
preferences: {}
users:
- name: {{.Username}}@{{.Tier}}
  user:
    auth-provider:
      config:
        client-id: {{ .ClientID }}
        client-secret: {{ .ClientSecret }}
        extra-scopes: groups
        id-token: {{.Token}}
        idp-issuer-url: {{ .Issuer }}
        refresh-token: {{.RefreshToken}}
      name: oidc
`

type configData struct {
	Cluster      string
	Clusters     []Clusters
	userName     string
	NS           string
	CAServerName string
	Token        string
	RefreshToken string
	Username     string
	ClientID     string
	Issuer       string
	ClientSecret string
	Tier         string
}

var funcs = template.FuncMap{"getCert": GetCertificate}

// NewKubeConfig returns an initialized KubeConfig struct.
func NewKubeConfig(cluster string, clusters []Clusters, username string, namespace string, caservername string, output io.ReadWriteCloser, clientID string, issuer string, clientSecret string) (*KubeConfig, error) {
	tmpl := template.Must(template.New("config").Funcs(funcs).Parse(content))

	log.Debug("Operating on clusters: ", clusters)
	log.Debug("Template: ", tmpl)

	if len(clusters) == 0 {
		log.Fatal("No clusters to operate on, check your config file! (hint: run with --debug)")
	}

	return &KubeConfig{
		cluster,
		clusters,
		username,
		namespace,
		caservername,
		tmpl,
		output,
		clientID,
		issuer,
		clientSecret,
		tier,
	}, nil
}

// Generate executes the writing of the config to the appropriate location (os.Stdout, os.File, etc...).
func (k *KubeConfig) Generate(token string, refreshToken string) error {
	defer k.Output.Close()

	log.Debug("kubeconfig struct contains: ", k)

	err := k.tmpl.Execute(k.Output, configData{
		k.Cluster,
		k.Clusters,
		k.Username,
		k.NS,
		k.CAServerName,
		token,
		refreshToken,
		userName,
		k.ClientID,
		k.Issuer,
		k.ClientSecret,
		k.Tier,
	})
	if err != nil {
		return err
	}

	return nil
}
