package fetch

import (
	"net/http"

	"bytes"

	"errors"

	"log"

	"github.com/morrah77/twitter-viewer/auth"
)

const ERR_MSG_NO_BYTES_READ = `No bytes were readfrom response`

type Fetcher interface {
	Fetch(string, string) (*bytes.Buffer, error)
}

type FetcherImpl struct {
	ApiUrl                 string
	Authorizer             auth.Authorizer
	Logger                 *log.Logger
	base64EncodedAuthToken string
}

func NewFetcher(apiUrl string, authorizer auth.Authorizer, logger *log.Logger) *FetcherImpl {
	return &FetcherImpl{
		ApiUrl:     apiUrl,
		Authorizer: authorizer,
		Logger:     logger,
	}
}

func (f *FetcherImpl) Fetch(screenName string, count string) (*bytes.Buffer, error) {
	var (
		err  error
		resp *http.Response
	)
	if f.base64EncodedAuthToken == `` {
		err = f.retryAuth(1)
		if err != nil {
			f.Logger.Println(`Fetch`, `Authorization error`, err)
			return nil, err
		}
	}
	req, err := http.NewRequest("GET", f.ApiUrl+`?screen_name=`+screenName+`&count=`+count, nil)
	if err != nil {
		f.Logger.Println(`Fetch`, `Could not create request`, err)
		return nil, err
	}
	req.Header.Set("Authorization", `Bearer `+f.base64EncodedAuthToken)
	req.Header.Set("Content-Type", `application/json`)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		f.Logger.Println(`Fetch`, `Could not get response`, err)
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		err = f.retryAuth(1)
		if err != nil {
			f.Logger.Println(`Fetch`, `Could not authorize`, err)
			return nil, err
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			f.Logger.Println(`Fetch`, `Could not get response`, err)
			return nil, err
		}
	}
	b := bytes.NewBuffer(make([]byte, 0))
	n, err := b.ReadFrom(resp.Body)
	if err != nil {
		f.Logger.Println(`Fetch`, `An error occured during response reading`, err)
		return b, err
	}
	if n <= 0 {
		f.Logger.Println(`Fetch`, ERR_MSG_NO_BYTES_READ)
		return b, errors.New(ERR_MSG_NO_BYTES_READ)
	}
	return b, nil
}

func (f *FetcherImpl) retryAuth(attemptsCount int) error {
	var err error
	for i := 0; i < attemptsCount; i++ {
		f.base64EncodedAuthToken, err = f.Authorizer.Authorize()
		if err == nil {
			f.Logger.Println(`Fetch`, `retryAuth success base64EncodedAuthToken`, f.base64EncodedAuthToken)
			return nil
		}
	}
	f.Logger.Println(`Fetch`, `retryAuth error base64EncodedAuthToken`, f.base64EncodedAuthToken)
	return err
}
