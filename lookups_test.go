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
func TestLookupRetrieve(t *testing.T) {
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

	lookup_search := "| rest /services/authentication/users splunk_server=local | table title realname roles"
	res, err := c.LookupRetrieve(lookup_search)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)

	kvlookup_search := "| rest /services/authentication/users splunk_server=local | rename title as key realname as value | table key value"
	kvres, err := c.LookupRetrieveKv(kvlookup_search)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", kvres)
}
