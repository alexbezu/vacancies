The main language is Golang
The service is called vacancies.

Functionality of the service:
 1. Stores WebHook and other configuration in environment variables.
 2. Stores new URLs in a Firestore instance.
 3. Sends new URLs via WebHook to a Telegram channel.
 4. A Cloud Run function is triggered by Google Scheduler twice a day. This function has a GET endpoint that starts the process of checking for new URLs.
 5. The algorithm for checking new URLs is as follows:
    a. Fetches a list of URLs from storage (Firestore).
    b. For each URL, performs a GET request and collects all urls from the page using `golang.org/x/net/html`.
    c. Compares the collected URLs with the existing URLs in the storage.
    d. Stores the new URLs in Firestore.
    e. Sends a notification for each new URL via WebHook.
 6. The project follows clean architecture principles.
 7. All dependencies of the service are implemented via interfaces (dependency injection).
 8. There is a local development server that can be used to trigger the URL processing manually.
 9. There are to types of urls: manually added to firestore (let's name them Sites) that used as page of list with vacancies (returned by GetSites function) and urls of actual vacancies (let's name them urls) for StoreURLs and GetURLs functions.

## Dependencies:
- `cloud.google.com/go/firestore` for Firestore integration.
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` for Telegram bot integration.
- `github.com/kelseyhightower/envconfig` for configuration management from environment variables.
- `github.com/sirupsen/logrus` for logging.
- `golang.org/x/net/html` for parsing HTML.

## Logging:
The service uses `logrus` for structured logging.

## Configuration:
The service uses `github.com/kelseyhightower/envconfig` to manage configuration from environment variables.
