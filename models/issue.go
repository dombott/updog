package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/prometheus/alertmanager/template"
)

type Issue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func (i *Issue) Hash() string {
	jsonBytes, _ := json.Marshal(i)
	return fmt.Sprintf("%x", sha256.Sum256(jsonBytes))
}

func IssueFromAlert(a template.Alert) *Issue {
	// TODO
	return &Issue{}
}
