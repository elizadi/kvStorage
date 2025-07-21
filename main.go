package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	handlers "kvStorage/internal/handler/router"
	store "kvStorage/internal/repository/storage"
	"kvStorage/internal/usecase"
	"os"
)

func main() {
	logger := logrus.New()

	logger.Info("Starting kvStorage")

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file. Error: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	storage, err := store.New(dbUrl, username, password, logger)
	if err != nil {
		logger.Fatalf("Failed start storage. Error: %v", err)
	}

	useCase := usecase.New(storage, logger)

	server := handlers.New(useCase, logger)

	router := handlers.Router(server)

	if err = router.Run(":8081"); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
