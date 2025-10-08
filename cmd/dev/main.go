package main

import (
	"log"
	"net/http"

	"github.com/alexbezu/vacancies/internal/config"
	"github.com/alexbezu/vacancies/internal/service"
	"github.com/alexbezu/vacancies/internal/storage"
	"github.com/alexbezu/vacancies/internal/webhook"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	_, err := config.FromEnv()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	// For now, we'll use in-memory storage and a logger-based webhook.
	// In the future, we can swap these out with other implementations.
	storage := storage.NewInMemoryStorage()
	webhook := webhook.NewLogWebhook(logger)

	svc := service.New(storage, webhook, logger)

	http.HandleFunc("/", checkNewURLsHandler(svc, logger))

	logger.Infof("listening on default port")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func checkNewURLsHandler(svc *service.Service, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := svc.ProcessURLs(r.Context()); err != nil {
			log.WithError(err).Error("failed to process URLs")
			http.Error(w, "failed to process URLs", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
