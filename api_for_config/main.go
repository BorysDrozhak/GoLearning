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

const config_file_name = "config.yml"

type StructConfig struct {
	A string `yaml:"a"`
	B string `yaml:"b"`
}

func main() {
	var configStruct StructConfig
	pwd := "/Users/bdrozhak/IdeaProjects/GoLearning/src/github.com/BorysDrozhak/GoLearning/api_for_config/"
	config_file_name_full_path := filepath.Join(pwd, config_file_name)

	r := gin.Default()
	v0 := r.Group("/v0")

	v0.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v0.GET("/read", func(c *gin.Context) {
		dat, err := ioutil.ReadFile(config_file_name_full_path)
		if err != nil {
			log.Println("Can't read file: ", err)
			c.Status(http.StatusInternalServerError)
		} else {
			err := yaml.Unmarshal(dat, &configStruct)
			if err != nil {
				log.Println("Can't read file: ", err)
				c.Status(http.StatusInternalServerError)
			} else {
				c.JSON(200, gin.H{
					"a": fmt.Sprintf("%v", configStruct.A),
					"b": fmt.Sprintf("%v", configStruct.B),
				})
			}
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
