package server

import (
    "io"
    "io/ioutil"
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
   "strconv"
)

const BUCKET_SLOTS = "SLOTS"
const KEY_LAST_SLOT = "last_slot"

type JSON map[string]interface{}

type Server struct {
    Host           string
    Port           string
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
	Config         Config
	Store          Store
}

type Config interface {
    GetServer() string
    GetLogin() string
    GetPassword() string
    GetCheckStatusUrl() string
    GetSendSmsUrl() string
    GetActiveSlots() []int
}

type Store interface {
    Get(bucket, key string) string
    Set(bucket, key, value string)
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

    lastSlot := getActiveSlot(s)

    params := make(map[string]string)
    params["u"] = s.Config.GetLogin()
    params["p"] = s.Config.GetPassword()
    params["l"] = lastSlot
    params["n"] = value
    params["n"] = "text"

    checkStatusUrl := s.Config.GetSendSmsUrl()
    resp, err := http.Post(checkStatusUrl, "application/json", nil)
    if err != nil {
        log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    log.Fatalln(err)
    }

    lgr.Printf(string(body))
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

func getActiveSlot(s Server) string {
    lastSlot, _ := strconv.ParseInt(s.Store.Get(BUCKET_SLOTS, KEY_LAST_SLOT), 10, 0)
    activeSlots := s.Config.GetActiveSlots()

    //lenSlots := len(activeSlots)

    var currentSlot int

    if activeSlots[lastSlot + 1] == 1 {
        currentSlot = int(lastSlot) + 1
    } else {
        currentSlot = 0
    }
    currentSlotStr := strconv.Itoa(currentSlot)
    s.Store.Set(BUCKET_SLOTS, KEY_LAST_SLOT, currentSlotStr)

    for _, slot := range activeSlots {
           lgr.Printf("%s\n", slot)
    }
    return currentSlotStr
}

func (s Server) checkStatus(w http.ResponseWriter, r *http.Request) {
     checkStatusUrl := s.Config.GetCheckStatusUrl()
     resp, err := http.Get(checkStatusUrl)
     if err != nil {
        log.Fatalln(err)
     }

      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        log.Fatalln(err)
     }

     lgr.Printf(string(body))

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
