#!/bin/bash

migrate -database "postgresql://$DBMASTERUSER:$DBMASTERPASS@tcp(database:5432)/$DBMASTERNAME?multiStatements=true" -path infrastucture/db/migrations $@