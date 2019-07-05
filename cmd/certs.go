package cmd

import (
	"crypto/tls"
	b64 "encoding/base64"
	"github.com/hashicorp/go-cleanhttp"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	//"net/http"
)

func GetCertificate(serverName string, address string) string {

	c := cleanhttp.DefaultClient()
	t := cleanhttp.DefaultTransport()
	t.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure,
		ServerName:         serverName,
	}
	c.Transport = t
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if address == "" {
		log.Fatal("No address supplied, cannot continue")
	}
	result, err := c.Get(address)
	if err != nil {
		log.WithFields(log.Fields{"server": address}).Fatal("Error getting remote certificate: ", err)
	}

	defer result.Body.Close()
	cert, err := ioutil.ReadAll(result.Body)

	encoded := b64.StdEncoding.EncodeToString(cert)

	return encoded

}
