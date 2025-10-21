# Vacancies

Vacancies is a Go service that monitors job websites for new vacancies and sends notifications through a webhook.

## Description

This service automates the process of checking for new job vacancies. It fetches a list of websites, scrapes them for job URLs, compares them with previously stored URLs, and sends notifications for new vacancies. The service is designed for deployment on Google Cloud Functions and is triggered by Google Scheduler.

## Features

-   Scrapes websites for URLs.
-   Stores URLs in a Firestore database.
-   Sends notifications for new URLs to a messenger via a bot.
-   Designed for serverless deployment on Google Cloud Functions.
-   Follows clean architecture principles.
-   Dependencies are managed using interfaces (dependency injection).
-   Includes a local development server.

## Architecture

The project follows the principles of Clean Architecture. The core logic is located in the `internal/service` package, which is independent of external frameworks and libraries. Dependencies such as storage and webhooks are injected into the service using interfaces, allowing for easy replacement of implementations.

-   `cmd`: Contains the entry points for the application (local development server and Cloud Function).
-   `internal`: Contains the core business logic of the application.
    -   `config`: Manages application configuration from environment variables.
    -   `service`: Implements the main functionality of the service.
    -   `storage`: Provides an interface for URL storage and its implementations.
    -   `webhook`: Provides an interface for sending notifications and its implementations.
-   `pkg`: Contains packages that can be shared with other applications.
    -   `bot`: A client for interacting with messengers.

## Getting Started

### Prerequisites

-   Go 1.19 or higher
-   Google Cloud SDK (for GCP deployment)

### Configuration

The service is configured using environment variables. The following variables are available:

| Variable                | Description                                  | 
| ----------------------- | -------------------------------------------- | 
| `VACANCIES_CHAT_ID`     | Your chat ID for the bot.                    | 
| `VACANCIES_BOT_TOKEN`   | Your secret token from BotFather.            | 
| `VACANCIES_GCP_PROJECT` | The ID of your Google project.               | 

### Adding URLs

1. Create a `sites` collection in Firestore and add a document with the following fields:

| Field Name | Field type | Value
| ---------- | -------------------------------------------- | -------
| createdAt  | timestamp  | now
| updatedAt  | timestamp  | now
| filter     | string     | https://www\.globallogic\.com/ua/careers/\S+-irc\d+
| name       | string     | GlobalLogic
| status     | string     | ""
| url        | string     | https://www.globallogic.com/ua/career-search-page/?keywords=golang&experience=none&location=ukraine/

2. You can test your regex filter using `test/filter_test.go`.


### Running the application

#### Local Development

Import your GCP DB data to local instance:

```bash
# Create a backup
gcloud firestore export gs://myexports

# Download them
gsutil -m cp -r gs://myexports ~/your/exports/

# Start an emulator with import
firebase emulators:start --import ~/your/exports/myexports/[TAB] --only firestore
```

To start the local development server, first start the Firestore emulator. Then, set the `VACANCIES_GCP_PROJECT`, `FIRESTORE_EMULATOR_HOST`, `VACANCIES_BOT_TOKEN`, and `VACANCIES_CHAT_ID` environment variables. Finally, run the following command:

```bash
go run cmd/dev/main.go service
```

You can trigger the URL processing by sending a GET request to `localhost:80`.

To see all the links from a job site, run:
```bash
go run cmd/dev/main.go links "https://www.work.ua/jobs-golang/"
```

To see the filtered links from a job site, run:
```bash
go run cmd/dev/main.go links "https://www.work.ua/jobs-golang/" "https://www\.work\.ua/jobs/\d+/"
```

#### Google Cloud Platform Deployment

The service is designed to be deployed on Google Cloud Functions. You can deploy it using the `gcloud` command-line tool.

1.  Set up a new GCP project and enable the Cloud Functions API.
2.  Set the `VACANCIES_CHAT_ID` and `VACANCIES_BOT_TOKEN` as secret environment variables. You can update them in `cmd/gcp/env.sh`.
3.  Deploy the function:
    ```bash
    cmd/gcp/deploy.sh
    ```
4.  Set up a Google Scheduler job to trigger the function's HTTP endpoint at regular intervals.

## Usage

The service is triggered by a GET request to its main endpoint. When triggered, it performs the following steps:

1.  Fetches a list of sites to scrape from storage.
2.  For each site, it performs a GET request and collects all URLs from the page.
3.  It compares the collected URLs with the existing URLs in storage.
4.  New URLs are stored in Firestore.
5.  A notification is sent for each new URL via the configured webhook.


## TODO

1. Remove old vacancies (50 days)
2. Add worker pool