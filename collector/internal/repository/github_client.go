package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gitubrr/GoSymGym/collector/internal/domain"
)

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
	userAgent  string
}

func NewGitHubClient(token string, timeout int) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		baseURL:   "https://api.github.com",
		token:     token,
		userAgent: "GoSymGym-CLI/1.0",
	}
}

func (c *GitHubClient) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)

	log.Printf("Requesting GitHub API: %s", url) // ← добавить

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err) // ← добавить
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
		log.Printf("Using GitHub token") // ← добавить
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err) // ← добавить
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("GitHub API response status: %d", resp.StatusCode) // ← добавить

	// ... остальной код

	switch resp.StatusCode {
	case http.StatusOK:
		var githubResp struct {
			Name        string    `json:"name"`
			FullName    string    `json:"full_name"`
			Description string    `json:"description"`
			Stars       int       `json:"stargazers_count"`
			Forks       int       `json:"forks_count"`
			Issues      int       `json:"open_issues_count"`
			Language    string    `json:"language"`
			HTMLURL     string    `json:"html_url"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&githubResp); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}

		return &domain.Repository{
			Name:        githubResp.Name,
			FullName:    githubResp.FullName,
			Description: githubResp.Description,
			Stars:       githubResp.Stars,
			Forks:       githubResp.Forks,
			Issues:      githubResp.Issues,
			Language:    githubResp.Language,
			HTMLURL:     githubResp.HTMLURL,
			CreatedAt:   githubResp.CreatedAt,
			UpdatedAt:   githubResp.UpdatedAt,
		}, nil

	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusForbidden:
		return nil, ErrRateLimit
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}
