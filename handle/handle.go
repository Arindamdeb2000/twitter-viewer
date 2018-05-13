package handle

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"errors"

	"log"

	"github.com/morrah77/twitter-viewer/fetch"
)

const ERR_MSG_NO_BYTES_WRITTEN = `No bytes were written to response`

type Handler struct {
	StrDefaultItemsNumber string
	Fetcher               fetch.Fetcher
	Logger                *log.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Println(`Handle`, `Accepted`, r.URL.String())
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.Logger.Println(`Handle`, `Could not parse query`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var num string
	num = h.StrDefaultItemsNumber
	if params[`num`] != nil && len(params[`num`]) > 0 {
		_, err := strconv.Atoi(params[`num`][0])
		if err != nil {
			h.Logger.Println(`Handle`, `Invalid num query parameter value`)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		num = params[`num`][0]
	}
	var screenName string
	segs := strings.Split(r.URL.Path[len(`/view/`):], `/`)
	if len(segs) != 1 {
		h.Logger.Println(`Handle`, `Invalid screen_name path segment in request`)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	screenName = segs[0]
	buffer, err := h.Fetcher.Fetch(screenName, num)
	if err != nil {
		h.Logger.Println(`Handle`, `Could not fetch items`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = Respond(Transform(buffer.Bytes()), w)
	if err != nil {
		h.Logger.Println(`Handle`, `Could not respond`, err)
	}
}

func Respond(body []byte, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", `application/json`)
	n, err := w.Write(body)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New(ERR_MSG_NO_BYTES_WRITTEN)
	}
	return nil
}

func Transform(fetchedRespBody []byte) []byte {
	//REM(h.lazar) fetched content transformation may be implemented here
	return fetchedRespBody
}
