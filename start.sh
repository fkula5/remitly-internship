#!/bin/sh

if [ -z "$1" ]; then
  echo "Usage: ./start.sh <port>"
  exit 1
fi

export API_EXTERNAL_PORT=$1

echo "Starting Stock Market Simulator on port $API_EXTERNAL_PORT with High Availability (2 instances)..."

docker-compose up --build -d --scale app=2

echo "Application is available at http://localhost:$API_EXTERNAL_PORT"
