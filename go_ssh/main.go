package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	// TODO: use ssh-agent instead
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Println("Private key is encrypted")
		passPhrase := os.Getenv("PASSPHRASE")
		if passPhrase == "" {
			log.Println("PASSPHRASE is not set as env variable")
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
		log.Println("enter host:", h)
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
		var stdout, stderr bytes.Buffer
		session.Stdout = &stdout
		session.Stderr = &stderr
		if err := session.Run(command); err != nil {
			log.Println("Failed to run command at host: ", h, "\n", err.Error())
		}
		c <- stderr.String() + stdout.String()

	}()
	return
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%s", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var port int
	var dstHosts arrayFlags
	//var dstHostsString, command string
	var command string

	flag.IntVar(&port, "p", 22, "an int")
	flag.Var(&dstHosts, "hosts", `destination addresses
	Usage: -hosts=host:port
	Examples:
	-hosts=127.{1,2,3}               -> [ 127.1:22, 127.2:22, 127.3:22 ]
	-hosts=127.{1,2:2222,3}          -> [ 127.1:22, 127.2:2222, 127.3:22 ]
	-hosts=127.{1:22,2:22,4} -p 2222 -> [ 127.1:22, 127.2:22, 127.4:2222 ]`)
	flag.StringVar(&command, "c", "hostname", "bash command you want to execute on the destination hosts")
	flag.Parse()

	if len(dstHosts) == 0 {
		log.Fatal("Destination Host is not defined")
	} else {
		log.Println("dstHosts: ", dstHosts)
	}

	if command == "" {
		log.Fatal("command can't be empty")
	} else {
		log.Println("command: ", command)
	}

	log.Println("default port: ", port)

	conf := makeUserConf()

	//// for test purposes (avoid env variables)
	//command = "sleep 3; hostname"

	var workers []chan string

	for _, dstHost := range dstHosts {
		if !strings.Contains(dstHost, ":") {
			dstHost += ":" + strconv.Itoa(port)
		}
		workers = append(workers, execCommand(dstHost, command, conf))
	}
	for _, w := range workers {
		fmt.Println(<-w)
	}
}
