#!/bin/bash

CONNSTR=$(jq -r '.db_url' ~/.gatorconfig.json)
cd ./internal/database/sql/schema
goose postgres "$CONNSTR" up
