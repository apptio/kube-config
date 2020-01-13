package cmd

import (
	"crypto/tls"
	b64 "encoding/base64"
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-cleanhttp"
	"io/ioutil"
	//"net/http"
)

//GetCertificate fetches remote certificate from the given serverName and address.
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
		log.WithFields(log.Fields{"server": address}).Error("Error getting remote certificate: ", err)
		return ""
	}

	defer result.Body.Close()
	cert, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.WithFields(log.Fields{"server": address}).Error("Failed to read results:", err)
		return ""
	}
	encoded := b64.StdEncoding.EncodeToString(cert)

	return encoded

}
