#!/bin/bash

source ./app-env

go build -o nomadx && ./nomadx
