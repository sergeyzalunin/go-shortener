package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	js "github.com/sergeyzalunin/go-shortener/internal/serializer/json"
	ms "github.com/sergeyzalunin/go-shortener/internal/serializer/msgpack"
	"github.com/sergeyzalunin/go-shortener/internal/shortener"
	"go.uber.org/zap"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	log             *zap.Logger
	redirectService shortener.RedirectService
}

func NewHandler(redirectService shortener.RedirectService, log *zap.Logger) RedirectHandler {
	return &handler{log, redirectService}
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

func (h *handler) httpError(w http.ResponseWriter, statusCode int, err error) {
	http.Error(w, http.StatusText(statusCode), statusCode)

	if err != nil {
		h.log.Error(err.Error(), zap.Error(err))
	}
}

// Get of RedirectHandler redirects the found in service request
// with StatusMovedPermanently.
// StatusNotFound has been returned in case the code not found in
// a repository. In other hand server returns with StatusInternalServerError.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Is(err, shortener.ErrRedirectNotFound) {
			h.httpError(w, http.StatusNotFound, err)

			return
		}

		h.httpError(w, http.StatusInternalServerError, err)

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
		h.httpError(w, http.StatusInternalServerError, err)

		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		h.httpError(w, http.StatusInternalServerError, err)

		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Is(err, shortener.ErrRedirectInvalid) {
			h.httpError(w, http.StatusBadRequest, err)

			return
		}

		h.httpError(w, http.StatusInternalServerError, err)

		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		h.httpError(w, http.StatusInternalServerError, err)

		return
	}

	setupResponse(w, contentType, http.StatusCreated, responseBody)
}
