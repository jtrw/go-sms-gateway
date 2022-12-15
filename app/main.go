package main
import (
  "log"
   "net/http"
   "github.com/pkg/errors"
   "github.com/go-chi/chi/v5"
   "github.com/go-chi/chi/v5/middleware"
   "github.com/jessevdk/go-flags"
   "github.com/jtrw/go-rest"

   server "sms-gateway/app/server"
)

type Server struct {
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
	Port           string
}

type Options struct {
    Port string `short:"p" long:"port" env:"SERVER_PORT" default:"8080" description:"Port web server"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    srv := server.Server {
        PinSize:   1,
        WebRoot:   "/",
        Version:   "1.0",
        Port: opts.Port,
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}




