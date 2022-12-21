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
func BenchmarkSearch(b *testing.B) {
	b.Logf("Openning connection with\nSPLUNK_HOST: %s\nSPLUNK_ADMIN_USERNAME: %s\nSPLUNK_ADMIN_PASSWORD: %s",
		os.Getenv("SPLUNK_HOST"), os.Getenv("SPLUNK_ADMIN_USERNAME"), os.Getenv("SPLUNK_ADMIN_PASSWORD"))
	c, err := gosplunk.NewClient(gosplunk.ClientConfig{
		Host:          os.Getenv("SPLUNK_HOST"),
		Username:      os.Getenv("SPLUNK_ADMIN_USERNAME"),
		Password:      os.Getenv("SPLUNK_ADMIN_PASSWORD"),
		App:           os.Getenv("SPLUNK_APP"),
		EnableLogging: true,
		Logger:        gosplunk.SimpleLogger{},
	})
	if err != nil {
		b.Fatal(err)
	}

	b.Log("Creating searches")

	for i := 0; i < b.N; i++ {
		_, err := c.SearchExport(gosplunk.NewSearch{
			Search:   "search index=_internal | head 5 | table _time",
			Earliest: "1",
			Latest:   "now",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
