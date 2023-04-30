package models

import (
	"strings"

	"github.com/prometheus/alertmanager/template"
)

type Issue struct {
	Title       string
	Body        string
	Labels      []string
	Fingerprint string
}

func IssueFromAlert(a template.Alert) *Issue {
	labels := strings.Split(a.Labels["updog/labels"], ";")
	labels = append(labels, "type/incident")
	return &Issue{
		Title:       a.Labels["updog/title"],
		Body:        a.Labels["updog/body"],
		Labels:      labels,
		Fingerprint: a.Fingerprint,
	}
}
