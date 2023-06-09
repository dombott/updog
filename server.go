package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dombott/updog/models"

	"github.com/google/go-github/v39/github"
	"github.com/prometheus/alertmanager/template"
)

const searchFormat = "updog:%s"
const titleFormat = "%s updog:%s"

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func (u *updog) webhook(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var data template.Data
	if err := dec.Decode(&data); err != nil {
		u.log.Error(err, "error decoding message")
		http.Error(w, "invalid request body", 400)
		return
	}

	for _, alert := range data.Alerts {
		if alert.Status == "firing" {
			if err := u.createIssueIfNotExists(*models.IssueFromAlert(alert)); err != nil {
				u.log.Error(err, "error creating issue")
				http.Error(w, "failed to create issue", 500)
				return
			}
		} else {
			if err := u.closeIssueIfExists(*models.IssueFromAlert(alert)); err != nil {
				u.log.Error(err, "error closing issue")
				http.Error(w, "failed to close issue", 500)
				return
			}
		}
	}
}

func (u *updog) createIssueIfNotExists(issue models.Issue) error {
	// search issue with identifier
	foundIssue, err := u.client.SearchIssue(fmt.Sprintf(searchFormat, issue.Fingerprint))
	if err != nil {
		return err
	}

	// issue with identifier already exists, do nothing
	if foundIssue != nil {
		u.log.Info("issue already exists, skipping", "fingerprint", issue.Fingerprint)
		return nil
	}

	labels, err := u.client.ListLabels()
	if err != nil {
		return err
	}

	// ensure labels exist
	for _, label := range issue.Labels {
		if !labelsContain(labels, label) {
			_, err := u.client.CreateLabel(label)
			if err != nil {
				return err
			}
		}
	}

	// create issue
	_, err = u.client.CreateIssue(fmt.Sprintf(titleFormat, issue.Title, issue.Fingerprint), issue.Body, issue.Labels)
	if err != nil {
		return err
	}
	return nil
}

func (u *updog) closeIssueIfExists(issue models.Issue) error {
	// search issue with identifier
	foundIssue, err := u.client.SearchIssue(issue.Fingerprint)
	if err != nil {
		return err
	}

	if foundIssue == nil {
		u.log.Info("no open issue found, skipping", "fingerprint", issue.Fingerprint)
		return nil
	}

	// issue with identifier exists, close
	if err := u.client.CloseIssue(*foundIssue.Number); err != nil {
		return err
	}

	return nil
}

func labelsContain(labels []*github.Label, lbl string) bool {
	for _, label := range labels {
		if *label.Name == lbl {
			return true
		}
	}
	return false
}
