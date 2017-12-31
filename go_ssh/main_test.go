package main

import (
	"testing"
	"io/ioutil"
	"os"
	"time"
)

var conf *config

func init() {
	InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	conf = makeUserConf()
}

func Test_execCommand(t *testing.T) {

	var workers = make(chan chan string)
	execCommand("127.0.0.1:22", "echo 1", conf, workers)

	expected := `--------------
host: 127.0.0.1:22
1
`
	select {
	case worker := <-workers:
		select {
		case result := <-worker:
			if result != expected {
				t.Error("Unexpected command result:", result)
			}
		case <-time.After(time.Second * 20):
			t.Error("Can't read output from chan because of timeout err")
			}
	case <-time.After(time.Second *2):
		t.Error("Can't get worker by execCommand func because of timeout err")

	}

}

//func Test_execCommandNegative(t *testing.T) {
//	var workers = make(chan chan string)
//	execCommand("127.0.0.1:222", "echo 1", conf, workers)
//	if result != expected {
//		t.Error("Unexpected command result:", result)
//	}
//}
