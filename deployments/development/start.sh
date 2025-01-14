#!/bin/bash

go mod tidy
go mod vendor

deployments/development/db-migration.sh up
exec air -c deployments/development/.air.toml