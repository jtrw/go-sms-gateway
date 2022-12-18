package main
import (
  "log"
    "fmt"
    "os"
   //"github.com/pkg/errors"
   "github.com/jessevdk/go-flags"
   "gopkg.in/yaml.v3"
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
    Config string `short:"f" long:"file" env:"CONF" description:"config file"`
}

type Config struct {
    Gateway struct {
         Server string `yaml:"server"`
         Login string `yaml:"login"`
         Password string `yaml:"password"`
    } `yaml:"gateway"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    config, errYaml := LoadConfig(opts.Config)
    if errYaml != nil {
        log.Println(errYaml)
    }
    fmt.Println(config)

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

func LoadConfig(file string) (*Config, error) {
	fh, err := os.Open(file) //nolint
	if err != nil {
		return nil, fmt.Errorf("can't load config file %s: %w", file, err)
	}
	defer fh.Close() //nolint

	res := Config{}
	if err := yaml.NewDecoder(fh).Decode(&res); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}
	return &res, nil
}




