#!/usr/bin/env bash
# This script is used to create migration the database schema.

migrate --version

migrate create -ext sql -dir ./migrations $1