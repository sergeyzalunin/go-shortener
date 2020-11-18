package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	js "github.com/sergeyzalunin/go-shortener/serializer/json"
	ms "github.com/sergeyzalunin/go-shortener/serializer/msgpack"
	"github.com/sergeyzalunin/go-shortener/shortener"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService}
}

func setupResponse(w http.ResponseWriter, contentType string, statusCode int, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.Redirect{}
	}

	return &js.Redirect{}
}

// Get of RedirectHandler redirects the found in service request
// with StatusMovedPermanently.
// StatusNotFound has been returned in case the code not found in
// a repository. In other hand server returns with StatusInternalServerError.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound))
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

// Post should reseive the Redirect struct representation inside the body.
// StatusInternalServerError code responses on each error
// StatusCreated code puts in response if all is ok and Redirect created.
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest))
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	setupResponse(w, contentType, http.StatusCreated, responseBody)
}
