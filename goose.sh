#!/bin/bash

cd sql/schema
goose postgres "$(cat ../../conn_string.txt)" "$1"
cd -
