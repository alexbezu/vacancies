#!/bin/bash
. ./cmd/gcp/env.sh

gcloud functions deploy CheckNewURLs \
--project $VACANCIES_GCP_PROJECT \
--runtime go125 \
--trigger-http \
--entry-point CheckNewURLs \
--region $VACANCIES_GCP_PROJECT_LOCATION \
--set-env-vars WEBHOOK=$VACANCIES_WEBHOOK,FIREBASE_PROJECT_ID=$VACANCIES_FIREBASE_PROJECT_ID \
--max-instances 1
