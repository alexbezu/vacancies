package vacancies

import (
	"context"
	"net/http"

	"github.com/alexbezu/vacancies/internal/config"
	"github.com/alexbezu/vacancies/internal/service"
	"github.com/alexbezu/vacancies/internal/storage"
	"github.com/alexbezu/vacancies/internal/webhook"

	"github.com/sirupsen/logrus"
)

var (
	svc *service.Service
	log *logrus.Logger
)

func init() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}

	storage, err := storage.NewFireStore(context.Background(), cfg.ProjectID)
	if err != nil {
		log.WithError(err).Fatal("failed to load db")
	}
	webhook := webhook.NewLogWebhook(log)

	svc = service.New(storage, webhook, log)
}

func CheckNewURLs(w http.ResponseWriter, r *http.Request) {
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
