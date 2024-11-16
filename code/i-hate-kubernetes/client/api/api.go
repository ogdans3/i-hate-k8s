package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/api/webhooks"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
)

func StartApiServer(ctx context.Context, port string, webhookChannel chan webhooks.WebhookPayload) {
	mux := http.NewServeMux()

	mux.HandleFunc("/webhook/github", func(w http.ResponseWriter, r *http.Request) {
		webhooks.GithubWebhookHandler(w, r, webhookChannel)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		console.InfoLog.Info("Starting API server on port ", port)
		if err := server.ListenAndServe(); err != nil {
			console.InfoLog.Fatal("Error starting server: ", err)
		}
	}()
	<-ctx.Done()
	console.InfoLog.Info("Shutting down API server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		console.InfoLog.Fatal("Error during server shutdown: ", err)
	}

	console.InfoLog.Info("API server stopped.")
}
