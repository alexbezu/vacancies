#!/bin/bash
. ./cmd/gcp/env.sh

gcloud functions deploy CheckNewURLs \
--project $VACANCIES_GCP_PROJECT \
--runtime go125 \
--trigger-http \
--entry-point CheckNewURLs \
--region $VACANCIES_GCP_PROJECT_LOCATION \
--set-env-vars VACANCIES_BOT_TOKEN=$VACANCIES_BOT_TOKEN,VACANCIES_CHAT_ID=$VACANCIES_CHAT_ID,VACANCIES_GCP_PROJECT=$VACANCIES_GCP_PROJECT, \
--max-instances 1
