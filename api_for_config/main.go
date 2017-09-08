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

type StructConfig struct {
	aggregation string `yaml:"aggregation"`
	b int `yaml:"chatTimeout"`
}


func main() {
	var configStruct StructConfig
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
	v0.GET("/config/:rule", func(c *gin.Context) {
		bytes, err := ioutil.ReadFile(configFileNameFullPath)
		if err != nil {
			log.Println("Can't read file: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		//log.Println(string(bytes))

		err = yaml.Unmarshal(bytes, &configStruct)
		if err != nil {
			log.Println("Can't read file: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		//fmt.Println(configStruct)
		switch rule := c.Params.ByName("rule"); rule {
			case "aggregation":
				v, err := yaml.Marshal(configStruct.aggregation)
				if err != nil {
					log.Println("Can't serialize: ", err)
					c.Status(http.StatusInternalServerError)
					return
				}
				c.YAML(200, v)
			case "chatTimeout":
				b := fmt.Sprintf("%v", configStruct.b)
				//fmt.Println(b)
				c.JSON(200, b)
			case "all":
				c.JSON(200, fmt.Sprintf("%v", &configStruct))
			default:
				//data, err  := json.Marshal(&configStruct)
				//if err != nil {
				//	log.Println("Can't load to json: ", err)
				//	c.Status(http.StatusInternalServerError)
				//	return
				log.Println("there is no such a rule: ", rule )
				c.Status(http.StatusNotFound)
				return
			}
		})
	r.Run(serverListenAddress) // listen and serve
}
