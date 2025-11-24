package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexbezu/vacancies/internal/config"
	"github.com/alexbezu/vacancies/internal/service"
	"github.com/alexbezu/vacancies/internal/storage"
	"github.com/alexbezu/vacancies/pkg/bot"
	"github.com/alexbezu/vacancies/pkg/scraper"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main service or ./main links 'https://www...' filter ")
		return
	}

	switch os.Args[1] {
	case "service":
		cfg, err := config.FromEnv()
		if err != nil {
			logger.WithError(err).Fatal("failed to load config")
		}

		scraper := scraper.NewSTD(logger)

		storage, err := storage.NewFireStore(context.Background(), cfg.ProjectID)
		if err != nil {
			logger.WithError(err).Fatal("failed to load db")
		}

		bot, err := bot.NewTelegram(cfg.BotToken, cfg.ChatID)
		if err != nil {
			logger.WithError(err).Error("failed to create a bot")
		}

		svc := service.New(storage, scraper, bot, logger)

		http.HandleFunc("/", checkNewURLsHandler(svc, logger))

		logger.Infof("listening on default port")
		if err := http.ListenAndServe(":80", nil); err != nil {
			log.Fatal(err)
		}
	case "links":
		url := os.Args[2]
		filter := ""
		if len(os.Args) > 3 {
			filter = os.Args[3]
		}
		scraper := scraper.NewColly(logger)
		links, _ := scraper.UrlsFromSite(context.TODO(), url, filter)
		for _, link := range links {
			fmt.Println(link)
		}
	default:
		fmt.Println("Usage: ./main service or ./main links 'https://www...' filter ")
		return
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
