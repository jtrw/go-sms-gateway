package server

import (
    "io"
    "log"
   "net/http"
   "strings"
   "github.com/pkg/errors"
   "github.com/go-chi/chi/v5"
   "github.com/go-chi/chi/v5/middleware"
   "github.com/go-chi/render"
   "github.com/jtrw/go-rest"
   lgr "github.com/go-pkgz/lgr"
   "encoding/json"
)

type JSON map[string]interface{}

type Server struct {
    Host           string
    Port           string
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

func (s Server) Run() error {
    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Port: %s", s.Port)

	if err := http.ListenAndServe(":"+s.Port, s.routes()); err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}

	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

    router.Use(middleware.Logger)
    router.Use(rest.Ping)

    router.Route("/api/v1", func(r chi.Router) {
        r.Post("/send/sms", s.sendSms)
        r.Get("/check/status/*", s.checkStatus)
    })

	return router
}

func (s Server) sendSms(w http.ResponseWriter, r *http.Request) {
    b, err := io.ReadAll(r.Body)
    if err != nil {
        lgr.Printf("[ERROR] %s", err)
    }
    value := string(b)
    lgr.Printf("[INFO] %s", value)

    //uri := chi.URLParam(r, "*")

    dataJson := &JSON{}
    //dataType := "text"
    if isContentTypeJson(r) {
        errJsn := json.Unmarshal([]byte(value), dataJson)
        if errJsn != nil {
            lgr.Printf("ERROR Invalid json in Data");
            return
        }
       // dataType = "json"
    }

    render.Status(r, http.StatusCreated)
    render.JSON(w, r, JSON{"status": "ok"})
    return
}

func (s Server) checkStatus(w http.ResponseWriter, r *http.Request) {
     render.Status(r, http.StatusCreated)
     render.JSON(w, r, JSON{"status": "ok"})
     return
}

func isContentTypeJson(r *http.Request) bool {
    return r.Header.Get("Content-Type") == strings.ToLower("application/json")
}

func getKeyAndBucketByUrl(uri string) (string, string) {
    chunks := strings.Split(uri, "/")

    length := len(chunks)
    keyStore := chunks[length-1]
    bucket := strings.Join(chunks[:length-1], "/")

    return keyStore, bucket
}
