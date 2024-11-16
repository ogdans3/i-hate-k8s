package webhooks

type WebhookPayload struct {
	Ref        string
	Repository struct {
		HtmlUrl  string
		GitUrl   string
		SshUrl   string
		CloneUrl string
	}
}
