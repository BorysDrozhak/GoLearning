package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type config struct {
	User       string
	Signer     ssh.Signer
	ClientConf *ssh.ClientConfig
}

func makeSigner(p string) ssh.Signer {
	key, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("unable to read private key: %v", key)
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Println("Private key is encrypted")
		passPhrase := os.Getenv("PASSPHRASE")
		if passPhrase == "" {
			log.Println("PASSPHRASE is not set as env variable")
			// TODO: ask paraphrase in case
		}
		rawKey, err := ssh.ParseRawPrivateKeyWithPassphrase(key, []byte(passPhrase))
		signer, err = ssh.NewSignerFromKey(rawKey)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}
	}
	return signer
}

func prepareSshConfig(s ssh.Signer, u string) *ssh.ClientConfig {

	c := &ssh.ClientConfig{
		User: u,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(s),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return c
}

func makeUserConf() *config {
	var c config

	c.User = os.Getenv("USER")
	if c.User == "" {
		log.Fatal("User is not defined")
	}

	home := "/Users/" + c.User
	private_key_path := "/.ssh/id_rsa"
	p := home + private_key_path

	signer := makeSigner(p)
	c.ClientConf = prepareSshConfig(signer, c.User)
	return &c
}

func execCommand(h string, command string, conf *config) (c chan string) {
	c = make(chan string)
	go func() {
		fmt.Println("enter host:", h)
		client, err := ssh.Dial("tcp", h, conf.ClientConf)
		if err != nil {
			log.Fatal("Failed to dial: ", err)
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Failed to create session: ", err)
		}
		defer session.Close()

		// Once a Session is created, you can execute a single command on
		// the remote side using the Run method.
		// TODO: up to 20 hosts; run a command simultaneity .. .. .
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run(command); err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}

		c <- b.String()
	}()
	return
}

func main() {
	c := makeUserConf()

	dstHosts := strings.Split(os.Getenv("H"), " ")
	if len(dstHosts) == 0 {
		log.Fatal("Destination Host is not defined")
	}
	command := os.Getenv("C")
	if command == "" {
		command = "env"
	}
	command = "date; hostname; sleep 3 ; date"
	//dstHosts="127.0.0.1:22 m104.sjc.opendns.com:22 graphite013.sjc.opendns.com:22 graphite011.sjc.opendns.com:22"
	var workers []chan string

	for _, dstHost := range dstHosts {
		workers = append(workers, execCommand(dstHost, command, c))
	}
	for _, w := range workers {
		fmt.Println(<- w)
	}
}
