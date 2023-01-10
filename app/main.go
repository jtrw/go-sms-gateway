package main
import (
  "log"
    "fmt"
    "os"
   //"github.com/pkg/errors"
   "github.com/jessevdk/go-flags"
   "gopkg.in/yaml.v3"
   server "sms-gateway/app/server"
   jstore "sms-gateway/app/store"
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
    Config string `short:"f" long:"file" env:"CONF" description:"config file" default:"configs.yaml"`
    StoragePath string `short:"s" long:"storage_path" default:"/var/tmp/jtrw_sms_gateway.db" description:"Storage Path"`
}

type Config struct {
    Gateway struct {
         Server string `yaml:"server"`
         Login string `yaml:"login"`
         Password string `yaml:"password"`
         CheckStatus string `yaml:"check_status"`
         SendSmsUrl string `yaml:"send_sms"`
         ActiveSlots []int `yaml:"active_slots"`
    } `yaml:"gateway"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    store := jstore.Store {
        StorePath: opts.StoragePath,
    }

    store.JBolt = store.NewStore()

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
        Config: config,
        Store: store,
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

func (c Config) GetServer() string {
    return c.Gateway.Server
}

func (c Config) GetLogin() string {
    return c.Gateway.Login
}

func (c Config) GetPassword() string {
    return c.Gateway.Password
}

func (c Config) GetCheckStatusUrl() string {
    return c.Gateway.Server + c.Gateway.CheckStatus
}

func (c Config) GetSendSmsUrl() string {
    return c.Gateway.Server + c.Gateway.SendSmsUrl
}

func (c Config) GetActiveSlots() []int {
    return c.Gateway.ActiveSlots
}



