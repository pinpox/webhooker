package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host        string                `yaml:"Host"`
	Port        string                `yaml:"Port"`
	GlobalToken string                `yaml:"GlobalToken"`
	Hooks       map[string]HookConfig `yaml:"Hooks"`
}

type HookConfig struct {
	Command string `yaml:"Command"`
	Token   string `yaml:"Token"`
}

func (h HookConfig) Run(hookName string) {
	log.Println("Running", hookName)
	// TODO Maybe set gin context as environment var?
	// TODO Maybe print the command output?
	// TODO check for error, return stdout, log stderr

	if h.Command != "" {
		cmd := exec.Command("/bin/sh", "-c", h.Command)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func (h HookConfig) Authorized(token string) bool {
	isAuthed := token == h.Token || token == config.GlobalToken
	log.Println("Authorization valid:", isAuthed)
	return isAuthed
}

func ParseConfig(path string) Config {

	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to open", path, err)
	}

	c := &Config{}
	if err := yaml.Unmarshal(buf, c); err != nil {
		log.Fatal(err)
	}

	if c.GlobalToken == "" {
		c.GlobalToken = os.Getenv("HOOKER_TOKEN")
	}

	return *c

}

var config Config

func main() {
	router := gin.Default()

	config = ParseConfig(os.Getenv("HOOKER_CONFIG"))

	router.GET("/:hook", func(c *gin.Context) {

		hookName := c.Param("hook")
		hook, ok := config.Hooks[hookName]

		if ok {

			// Check token exists and is authorized
			tokenHeaders, tokenFound := c.Request.Header["Token"]
			if !tokenFound || !hook.Authorized(tokenHeaders[0]) {
				c.Status(http.StatusUnauthorized)
				return
			}

			// TODO: should we return status 200 before command is done for
			// long-running commands?
			// go val.Run(hookName)
			hook.Run(hookName)
			c.Status(http.StatusOK)

		} else {
			log.Println(hookName, "does not exist")
			c.Status(http.StatusNotFound)
		}

	})

	router.Run(config.Host + ":" + config.Port)
}
