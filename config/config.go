package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	// AppVersion Current application version
	AppVersion string `long:"app_version" env:"APP_VERSION" description:"Application version" required:"false" default:"v1.0.0"`
	// ListenAddr HTTP server address
	ListenAddr string `long:"port" env:"LISTEN_ADDR" description:"Listen to port (format: :3000|127.0.0.1:3000)" required:"false" default:":3000"`
	// BasePath HTTP server base URL
	BasePath string `long:"base_path" env:"BASE_PATH" description:"Base path of the host" required:"false" default:"/api/packer"`
}

func NewConfig() *Configuration {
	config := &Configuration{}
	config.parseConfig()
	return config
}

func (c *Configuration) PrintValues() {
	s := reflect.ValueOf(c).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%-25s: %70v\n", typeOfT.Field(i).Name, f.Interface())
	}
}

func (c *Configuration) parseConfig() {
	p := flags.NewParser(c, flags.Default)
	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] error parsing config:", err)
		var flagsErr *flags.Error
		if ok := errors.Is(err, flagsErr); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
