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

// func seed() error {
// 	log := logrus.New()

// 	cfg, err := config.FromEnv()
// 	if err != nil {
// 		log.WithError(err).Error("failed to load config")
// 		return err
// 	}

// 	ctx := context.Background()
// 	fs, err := firestore.NewClient(ctx, cfg.ProjectID)
// 	if err != nil {
// 		log.WithError(err).Error("failed to open firestore")
// 		return err
// 	}

// 	_, _, err = fs.Collection("sites").Add(ctx, model.JobSite{
// 		Name:          "Djinni",
// 		URL:           "https://djinni.co/jobs/keyword-golang/",
// 		CreatedAt:     time.Now(),
// 		ScraperConfig: map[string]any{"filter": ""},
// 		UpdatedAt:     time.Now(),
// 	})
// 	if err != nil {
// 		log.WithError(err).Error("failed to add a seed")
// 		return err
// 	}

// 	return nil
// }

// var _ = seed()

func init() {
	log := logrus.New()

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
