#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=kawajiritatsuyoshi -cache=false -production=false