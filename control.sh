#!/bin/sh

if [ -z "$1" ]; then
 echo 'Please specify a command [build|test|bench|run]'
fi

case "$1" in
  build)
    echo "Building..."
    rm -rf ./twitter-viewer && go build -o ./twitter-viewer && chmod u+x ./twitter-viewer
     if [ $? = 0 ]; then
       echo "Built successfully"
     else
       echo 'Build failed'
     fi ;;
  test)
    echo "Running tests..."
    go test ./auth/ && \
    go test ./fetch/ && \
    go test ./handle/ ;;
  bench)
    echo "Running tests with benchmarks..."
    go test ./auth/ -bench && \
    go test ./fetch/ -bench . && \
    go test ./handle/ -bench . ;;
  run)
    echo "Running..."
    ./twitter-viewer --listen-addr :8080 --api-key $2 --api-secret $3
esac
