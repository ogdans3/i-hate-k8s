package webhooks

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

// Payload represents the expected structure of incoming webhook data
type Payload struct {
	Ref        string `json:"ref"`
	Repository struct {
		HtmlUrl  string `json:"html_url"`
		GitUrl   string `json:"git_url"`
		SshUrl   string `json:"ssh_url"`
		CloneUrl string `json:"clone_url"`
	} `json:"repository"`
}

func GithubWebhookHandler(w http.ResponseWriter, r *http.Request, webhookChannel chan WebhookPayload) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	console.InfoLog.Info("Received webhook: ", payload.Ref, payload.Repository.GitUrl, "\n")
	webhookChannel <- WebhookPayload{
		Ref: payload.Ref,
		Repository: struct {
			HtmlUrl  string
			GitUrl   string
			SshUrl   string
			CloneUrl string
		}(payload.Repository),
	}
}
