package app

import (
	"encoding/base64"
	"os"

	"github.com/gorilla/securecookie"
	"gopkg.in/ini.v1"
)

type server struct {
	ListenAddress string `ini:"Listen"`
	Port          string `ini:"Port"`
}

type session struct {
	AuthenticationKey string `ini:"AuthenticationKey"`
	EncryptionKey     string `ini:"EncryptionKey"`
}

type database struct {
	File    string `ini:"File"`
	Install bool   `ini:"-"`
}

type logging struct {
	File string `ini:"File"`
}

type ConfigModel struct {
	Server   *server
	Session  *session
	Database *database
	Logging  *logging
}

func (a *app) InitConfig(file string) {

	cfg := &ConfigModel{

		Server: &server{
			ListenAddress: "",
			Port:          "8000",
		},

		Session: &session{
			AuthenticationKey: base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)),
			EncryptionKey:     base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)),
		},

		Database: &database{
			File:    "fpsmonitor.sqlite",
			Install: false,
		},

		Logging: &logging{
			File: "fpsmonitor.log",
		},
	}

	_, err := os.Stat(file)
	if err != nil {
		if !os.IsExist(err) {

			iniCfg := ini.Empty()
			secServer, _ := iniCfg.NewSection("Server")
			secServer.ReflectFrom(&cfg.Server)
			secSession, _ := iniCfg.NewSection("Session")
			secSession.ReflectFrom(&cfg.Session)
			secDatabase, _ := iniCfg.NewSection("Database")
			secDatabase.ReflectFrom(&cfg.Database)
			secLogging, _ := iniCfg.NewSection("Logging")
			secLogging.ReflectFrom(&cfg.Logging)
			iniCfg.SaveTo(file)

		}
	}

	iniCfg, _ := ini.Load(file)

	iniCfg.Section("Server").MapTo(&cfg.Server)
	iniCfg.Section("Session").MapTo(&cfg.Session)
	iniCfg.Section("Database").MapTo(&cfg.Database)
	iniCfg.Section("Logging").MapTo(&cfg.Logging)

	a.config = cfg
}

func (a *app) SetConfig(cfg *ConfigModel) {
	a.config = cfg
}

func (a *app) Config() *ConfigModel {
	if a.config == nil {
		a.config = &ConfigModel{

			Server: &server{
				ListenAddress: "",
				Port:          "8000",
			},

			Session: &session{
				AuthenticationKey: base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)),
				EncryptionKey:     base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)),
			},

			Database: &database{
				File:    "fpsmonitor.sqlite",
				Install: false,
			},

			Logging: &logging{
				File: "fpsmonitor.log",
			},
		}
	}

	return a.config
}
