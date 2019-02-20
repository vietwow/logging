#!/bin/bash

glide init --non-interactive

glide up --strip-vendor

go run main.go
