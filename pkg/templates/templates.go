package templates

import (
	"html/template"
	"log"
	"net/http"
)

type tokenTmplData struct {
	IDToken      string
	RefreshToken string
	RedirectURL  string
	Claims       string
}

var tokenTmpl = template.Must(template.New("token.html").Parse(`<html>
  <head>
    <style>
/* make pre wrap */
pre {
 white-space: pre-wrap;       /* css-3 */
 white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
 white-space: -pre-wrap;      /* Opera 4-6 */
 white-space: -o-pre-wrap;    /* Opera 7 */
 word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
    </style>
  </head>
  <body>
	  <p> Your token has been generated! </p>
    <p> Token: <pre><code>{{ .IDToken }}</code></pre></p>
		<p> The following information was provided by the authentication provider </p>
    <p> Claim Details: <pre><code>{{ .Claims }}</code></pre></p>
	{{ if .RefreshToken }}
    <p> Refresh Token: <pre><code>{{ .RefreshToken }}</code></pre></p>
	<form action="{{ .RedirectURL }}" method="post">
	  <input type="hidden" name="refresh_token" value="{{ .RefreshToken }}">
	  <input type="submit" value="Redeem refresh token">
    </form>
	{{ end }}
	  <p> Your Kubeconfig file has been generated - this page will remain open for your reference </p>
  </body>
</html>
`))

//RenderToken renders the token template in html with passed in values.
func RenderToken(w http.ResponseWriter, redirectURL, idToken, refreshToken string, claims []byte) {
	renderTemplate(w, tokenTmpl, tokenTmplData{
		IDToken:      idToken,
		RefreshToken: refreshToken,
		RedirectURL:  redirectURL,
		Claims:       string(claims),
	})
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *template.Error:
		// An ExecError guarantees that Execute has not written to the underlying reader.
		log.Printf("Error rendering template %s: %s", tmpl.Name(), err)

		// TODO(ericchiang): replace with better internal server error.
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		// An error with the underlying write, such as the connection being
		// dropped. Ignore for now.
	}
}
