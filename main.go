package main

import (
	"flag"

	"net/http"

	"strconv"

	"log"

	"os"

	"github.com/morrah77/twitter-viewer/auth"
	"github.com/morrah77/twitter-viewer/fetch"
	"github.com/morrah77/twitter-viewer/handle"
)

const VENDOR_API_URL = `https://api.twitter.com`

const ERROR_MSG_NO_API_KEY = `No API key passed!`
const ERROR_MSG_NO_API_SECRET = `No API secret passed!`
const ERROR_MSG_INCORRECT_DEFAULT_ITEMS_NUMBER = `Incorrect Default items number!`

type Conf struct {
	ApiKey             string
	ApiSecret          string
	ListenAddr         string
	DefaultItemsNumber int
}

var (
	conf   Conf
	logger *log.Logger
)

func init() {
	flag.StringVar(&conf.ApiKey, `api-key`, ``, `Twitter API key (Consumer Key)`)
	flag.StringVar(&conf.ApiSecret, `api-secret`, ``, `Twitter API secret (Consumer Secret)`)
	flag.StringVar(&conf.ListenAddr, `listen-addr`, `:8080`, `Address to listen`)
	flag.IntVar(&conf.DefaultItemsNumber, `default-items`, 10, `Default items number`)
	flag.Parse()
	logger = log.New(os.Stdout, `Twitter viewer`, log.Flags())
}

func main() {
	validateEnv()
	var (
		authorizer auth.Authorizer
		fetcher    fetch.Fetcher
		handler    http.Handler
		err        error
	)
	strDefaultItemsNumber := strconv.Itoa(conf.DefaultItemsNumber)
	authorizer = auth.NewAuthorizer(
		VENDOR_API_URL+`/oauth2/token`,
		conf.ApiKey,
		conf.ApiSecret,
		logger)
	fetcher = fetch.NewFetcher(
		VENDOR_API_URL+`/1.1/statuses/user_timeline.json`,
		authorizer,
		logger)
	handler = &handle.Handler{
		StrDefaultItemsNumber: strDefaultItemsNumber,
		Fetcher:               fetcher,
		Logger:                logger,
	}
	http.HandleFunc(`/view/`, handler.ServeHTTP)
	logger.Println(`Listen on`, conf.ListenAddr)
	err = http.ListenAndServe(conf.ListenAddr, nil)
	logger.Printf(`error %#v`, err)
	failOnError(err)
}

func validateEnv() {
	if conf.ApiKey == `` {
		panic(ERROR_MSG_NO_API_KEY)
	}
	if conf.ApiSecret == `` {
		panic(ERROR_MSG_NO_API_SECRET)
	}
	if conf.DefaultItemsNumber <= 0 {
		panic(ERROR_MSG_INCORRECT_DEFAULT_ITEMS_NUMBER)
	}
}

func failOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
