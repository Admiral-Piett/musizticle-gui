package config

import (
    "encoding/json"
    "github.com/Admiral-Piett/musizticle-gui/app/utils"
    "gitlab.com/avarf/getenvs"
    "log"
    "os"
    "reflect"
)

type Settings struct {
    ConfigFilePath     string `env:"CONFIG_FILE" default:"../../config"`
    BackendHost        string `env:"BACKEND_HOST" default:"http://localhost:9000"`
    BackendUsername    string
    BackendPassword    string
    BackendAuthToken   string
    BackendReauthToken string
}

// TODO - 2
type ConfigFile struct {
    BackendUsername string `json:"backendUsername"`
    BackendPassword string `json:"backendPassword"`
    //BackendReauthToken string `json:"backendReauthToken"`
}

func NewSettings() *Settings {
    s := &Settings{}
    s.load()
    return s
}

func (s *Settings) load() {
    keys := utils.GetStructFieldNames(*s)

    t := reflect.TypeOf(*s)
    for _, key := range keys {
        field, ok := t.FieldByName(key)
        if !ok {
            continue
        }
        envKey := field.Tag.Get("env")
        if envKey == "" {
            continue
        }
        value := getenvs.GetEnvString(envKey, field.Tag.Get("default"))
        s.set(key, value)
    }
    s.loadLocalFile()
}

func (s *Settings) loadLocalFile() {
    if s.ConfigFilePath == "" {
        return
    }
    data, err := os.ReadFile(s.ConfigFilePath)
    if err != nil {
        return
    }
    c := ConfigFile{}
    err = json.Unmarshal(data, &c)
    if err != nil {
        return
    }
    // TODO - 2
    s.BackendUsername = c.BackendUsername
    s.BackendPassword = c.BackendPassword
}

func (s *Settings) set(key string, value interface{}) {
    reflect.ValueOf(s).Elem().FieldByName(key).Set(reflect.ValueOf(value))
}

func (s *Settings) SetConfigFile() {
    c := ConfigFile{
        // TODO - 2
        BackendUsername: s.BackendUsername,
        BackendPassword: s.BackendPassword,
    }
    jsonBody, err := json.Marshal(c)
    if err != nil {
        log.Printf("SaveToFileFailure: %s\n", err)
    }
    _ = os.WriteFile(s.ConfigFilePath, jsonBody, 0777)
    log.Println("ConfigFile save complete")
}
