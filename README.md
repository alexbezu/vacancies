# Vacancies

Vacancies is a Go service that monitors job websites for new vacancies and sends notifications through a webhook.

## Description

This service is designed to automate the process of checking for new job vacancies. It fetches a list of websites, scrapes them for job URLs, compares them with already stored URLs, and sends notifications for new findings. The service is built to be deployed on Google Cloud Run and triggered by Google Scheduler.

## Features

-   Scrapes websites for URLs.
-   Stores URLs in a Firestore database.
-   Sends notifications for new URLs via webhooks.
-   Designed for serverless deployment on Google Cloud Run Functions.
-   Follows clean architecture principles.
-   Dependencies are managed through interfaces (dependency injection).
-   Includes a local development server.

## Architecture

The project follows the principles of Clean Architecture. The core logic is located in the `internal/service` package, which is independent of external frameworks and libraries. Dependencies such as storage and webhooks are injected into the service using interfaces, allowing for easy replacement of implementations.

-   `cmd`: Contains the entry points for the application (local development server and Cloud Run function).
-   `internal`: Contains the core business logic of the application.
    -   `config`: Manages application configuration from environment variables.
    -   `service`: Implements the main functionality of the service.
    -   `storage`: Provides an interface for URL storage and its implementations.
    -   `webhook`: Provides an interface for sending notifications and its implementations.
-   `pkg`: Contains packages that can be shared with other applications.
    -   `firebase`: A client for interacting with Firebase.

## Getting Started

### Prerequisites

-   Go 1.19 or higher
-   Google Cloud SDK (for GCP deployment)

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/alexbezu/vacancies.git
    cd vacancies
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```

### Configuration

The service is configured using environment variables. The following variables are available:

| Variable              | Description                                  | Default |
| --------------------- | -------------------------------------------- | ------- |
| `WEBHOOK`             | The URL of the webhook for notifications.    |         |
| `FIREBASE_PROJECT_ID` | The ID of your Firebase project.             |         |

### Running the application

#### Local Development

To start the local development server, run the following command:

```bash
go run cmd/dev/main.go
```

You can trigger the URL processing by sending a GET request to the `/` endpoint.

#### Google Cloud Platform Deployment

The service is designed to be deployed on Google Cloud Functions. You can deploy it using the `gcloud` command-line tool.

1.  Set up a new GCP project and enable the Cloud Functions API.
2.  Deploy the function:
    ```bash
    gcloud functions deploy CheckNewURLs --runtime go119 --trigger-http --entry-point CheckNewURLs --region <YOUR_REGION> --set-env-vars WEBHOOK=<YOUR_WEBHOOK_URL>,FIREBASE_PROJECT_ID=<YOUR_GCP_PROJECT_ID>
    ```
3.  The command will deploy the `CheckNewURLs` function.
4.  Set up a Google Scheduler job to trigger the function's HTTP endpoint at regular intervals.

## Usage

The service is triggered by a GET request to its main endpoint. When triggered, it performs the following steps:

1.  Fetches a list of sites to scrape from the storage.
2.  For each site, it performs a GET request and collects all the URLs from the page.
3.  It compares the collected URLs with the existing URLs in the storage.
4.  New URLs are stored in Firestore.
5.  A notification is sent for each new URL via the configured webhook.
