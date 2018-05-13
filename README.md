#Simple twitter viewer

Accepts HTTP GET request like

`/view/<twitter_account>?num=N`

where `<twitter_account>` is Twitter `screen_name`,

`N` is number of tweets to respond.

Responds with `N` last tweets from specified `<twitter_account>`.

##Usage example

`curl -iv http://localhost:8080/view/helenlazar?num=10`

##Build

###Manually

`go build -o twitter-viewer main.go`

or

`./control.sh build`

###With Docker

`docker build -t twitter-viewer -f Dockerfile .`


#Run

###Manually

`./twitter-viewer --listen-addr :8080 --api-key <your_twitter_app_API_key> --api-secret <your_twitter_app_API_secret>`

or

`./control.sh ./twitter-viewer <your_twitter_app_API_key> <your_twitter_app_API_secret>`

###With Docker

`docker run --rm -d --name twitter-viewer -p 8080:8080 twitter-viewer ./twitter-viewer  "--listen-addr=:8080" "--api-key=<your_twitter_app_API_key>" "--api-secret=<your_twitter_app_API_secret>"`

##Test

###Unit tests

```
go test ./auth/ && \

go test ./fetch/ && \
 
go test ./handle/
```

or

`./control.sh test`

###Manually

`curl -iv http://localhost:8080/view/helenlazar?num=10`

###With browser

`http://localhost:8080/view/helenlazar?num=10`
