package gosplunk_test

import (
	"os"
	"testing"

	gosplunk "github.com/kyzykyky/go-splunk"
)

// In order to test, you need to set the following environment variables:
// SPLUNK_HOST
// SPLUNK_ADMIN_USERNAME
// SPLUNK_ADMIN_PASSWORD
// SPLUNK_USER_USERNAME
// SPLUNK_USER_TOKEN
func TestNewClient(t *testing.T) {
	t.Logf("Openning connection with\nSPLUNK_HOST: %s\nSPLUNK_ADMIN_USERNAME: %s\nSPLUNK_ADMIN_PASSWORD: %s",
		os.Getenv("SPLUNK_HOST"), os.Getenv("SPLUNK_ADMIN_USERNAME"), os.Getenv("SPLUNK_ADMIN_PASSWORD"))
	c, err := gosplunk.NewClient(gosplunk.ClientConfig{
		Host:          os.Getenv("SPLUNK_HOST"),
		Username:      os.Getenv("SPLUNK_ADMIN_USERNAME"),
		Password:      os.Getenv("SPLUNK_ADMIN_PASSWORD"),
		EnableLogging: true,
		Logger:        gosplunk.SimpleLogger{},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", c)

	t.Logf("Creating subclient search with\nSPLUNK_USER_USERNAME: %s\nSPLUNK_USER_TOKEN: %s",
		os.Getenv("SPLUNK_USER_USERNAME"), os.Getenv("SPLUNK_USER_TOKEN"))
	_, err = c.SubClient(os.Getenv("SPLUNK_USER_USERNAME"), os.Getenv("SPLUNK_USER_TOKEN")).SearchExport(gosplunk.NewSearch{
		Search:   "search index=_internal | head 5 | table _time",
		Earliest: "1",
		Latest:   "now",
	})
	if err != nil {
		t.Fatal(err)
	}
}
