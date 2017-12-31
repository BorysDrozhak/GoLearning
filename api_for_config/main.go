package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"fmt"
)

const configFileName = "config.yml"
const serverListenAddress = "127.0.0.1:8090"

type StructConfigJson struct {
	Aggregation string `json:"aggregation"`
	ChatTimeout string `json:"chatTimeout"`
	URL string `json:"url"`
}

type StructConfig struct {
	Aggregation string `yaml:"aggregation"`
	ChatTimeout string `yaml:"chatTimeout"`
	Url string `yaml:"url"`
}

func main() {
	var ConfigStruct StructConfig
	pwd := "/Users/bdrozhak/IdeaProjects/GoLearning/src/github.com/BorysDrozhak/GoLearning/api_for_config/"
	configFileNameFullPath := filepath.Join(pwd, configFileName)

	r := gin.Default()
	v0 := r.Group("/v0")

	v0.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"version": "v0",
		})
	})
	v0.GET("/configs/:rule", func(c *gin.Context) {
		conf, err := ioutil.ReadFile(configFileNameFullPath)
		if err != nil {
			log.Println("Can't read file: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		err = yaml.Unmarshal(conf, &ConfigStruct)
		if err != nil {
			log.Println("Can't read file: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		switch rule := c.Params.ByName("rule"); rule {
			case "aggregation":
				c.YAML(200, ConfigStruct.Aggregation)
			case "chatTimeout":
				c.YAML(200, ConfigStruct.ChatTimeout)
			case "all":
				c.YAML(200, ConfigStruct)
			default:
				log.Println("there is no such a rule: ", rule )
				c.Status(http.StatusNotFound)
				return
			}
		})

	v0.POST("/url", func(c *gin.Context) {
		var configStructJson StructConfigJson
		c.BindJSON(&configStructJson)
		ConfigStruct.Url = fmt.Sprintf("%v",configStructJson)
		fmt.Println(ConfigStruct.Url)

	})

	r.Run(serverListenAddress) // listen and serve
}
