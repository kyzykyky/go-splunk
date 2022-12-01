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
func TestSavedSearch(t *testing.T) {
	t.Logf("Openning connection with\nSPLUNK_HOST: %s\nSPLUNK_ADMIN_USERNAME: %s\nSPLUNK_ADMIN_PASSWORD: %s",
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
		t.Fatal(err)
	}

	t.Log("Creating saved search")
	searchName := "gosplunk_test_saved_search"
	err = c.SavedSearchCreate(gosplunk.SavedSearch{
		Name:        searchName,
		Search:      "search index=_internal | head 5 | table _time",
		Description: "gosplunk test saved search",
	}, gosplunk.NameSpace{})
	if err == gosplunk.ErrConflict {
		t.Log("Test saved search already exists")
	} else if err != nil {
		t.Fatal(err)
	}

	t.Log("Updating saved search with schedule")
	cronSchedule := "*/5 * * * *"
	err = c.SavedSearchUpdate(gosplunk.SavedSearch{
		Name:              searchName,
		IsScheduledString: "true",
		Cron:              cronSchedule,
	}, gosplunk.NameSpace{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Getting saved search")
	ssearch, err := c.SavedSearchGet(searchName, gosplunk.NameSpace{})
	if err != nil {
		t.Fatal(err)
	}
	if ssearch.Cron != cronSchedule {
		t.Fatalf("Expected cron schedule %s, got %s", cronSchedule, ssearch.Cron)
	}

	t.Log("Deleting saved search")
	err = c.SavedSearchDelete(searchName, gosplunk.NameSpace{})
	if err != nil {
		t.Fatal(err)
	}
}
