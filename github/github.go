package github

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type Client struct {
	githubClient *github.Client
	owner        string
	repo         string

	labelCache []*github.Label
	cacheAge   time.Time
	cacheTTL   time.Duration
}

// NewClient creates a new GitHub API client using the provided OAuth2 token, owner and repo name.
func NewClient(token, owner, repo string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		githubClient: github.NewClient(tc),
		owner:        owner,
		repo:         repo,
		cacheTTL:     5 * time.Minute, // TODO make configurable
	}
}

// CreateIssue creates a new issue on the repository.
func (c *Client) CreateIssue(title, body string, labels []string) (*github.Issue, error) {
	issueRequest := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: &labels,
	}
	issue, _, err := c.githubClient.Issues.Create(context.Background(), c.owner, c.repo, issueRequest)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

// SearchIssue searches for existing issues where the identifier is included in the title
func (c *Client) SearchIssue(identifier string) (*github.Issue, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	query := fmt.Sprintf("repo:%s/%s %s in:title", c.owner, c.repo, identifier)
	results, _, err := c.githubClient.Search.Issues(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	if *results.Total > 0 {
		return results.Issues[0], nil
	}
	return nil, nil
}

// CloseIssue closes the specified issue on the repository.
func (c *Client) CloseIssue(number int) error {
	_, _, err := c.githubClient.Issues.Edit(context.Background(), c.owner, c.repo, number, &github.IssueRequest{State: github.String("closed")})
	if err != nil {
		return err
	}
	return nil
}

// ListLabels lists the labels on the repository.
func (c *Client) ListLabels() ([]*github.Label, error) {
	if time.Since(c.cacheAge) < time.Minute {
		return c.labelCache, nil
	}

	labels, _, err := c.githubClient.Issues.ListLabels(context.Background(), c.owner, c.repo, nil)
	if err != nil {
		return nil, err
	}

	c.labelCache = labels
	c.cacheAge = time.Now()

	return labels, nil
}

// CreateLabel creates a new label on the repository.
func (c *Client) CreateLabel(name string) (*github.Label, error) {
	labelRequest := &github.Label{Name: github.String(name), Color: github.String(randomColor())}
	label, _, err := c.githubClient.Issues.CreateLabel(context.Background(), c.owner, c.repo, labelRequest)
	if err != nil {
		return nil, err
	}

	// invalidate the label cache
	c.cacheAge = time.Time{}

	return label, nil
}

func randomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
