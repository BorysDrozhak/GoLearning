package main

import (
	//"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
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

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLog(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func makeSigner(p string) ssh.Signer {
	key, err := ioutil.ReadFile(p)
	if err != nil {
		Error.Fatalf("unable to read private key: %v", key)
	}
	// TODO: use ssh-agent instead
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		Info.Println("Private key is encrypted")
		passPhrase := os.Getenv("PASSPHRASE")
		if passPhrase == "" {
			Info.Println("PASSPHRASE is not set as env variable")
		}
		rawKey, err := ssh.ParseRawPrivateKeyWithPassphrase(key, []byte(passPhrase))
		signer, err = ssh.NewSignerFromKey(rawKey)
		if err != nil {
			Error.Fatalf("unable to parse private key: %v", err)
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
		Error.Fatal("User is not defined")
	}

	home := "/Users/" + c.User
	privateKeyPath := "/.ssh/id_rsa"
	p := home + privateKeyPath

	signer := makeSigner(p)
	c.ClientConf = prepareSshConfig(signer, c.User)
	return &c
}

func execCommand(h string, command string, conf *config, workers chan chan string) {
	go func() {
		c := make(chan string, 1)
		//workers <- c
		defer func() { workers <- c }()
		Trace.Println("enter host:", h)
		client, err := ssh.Dial("tcp", h, conf.ClientConf)
		if err != nil {
			Error.Println(
				"Failed to dial", h)
			c <- err.Error()
			return
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			Error.Println(
				"Failed to create session with: " + h)
			c <- err.Error()
			return
		}
		defer session.Close()

		if response, err := session.CombinedOutput(command); err != nil {
			c <- "--------------\n" + "host: " + h + "\n" + string(response) + "\nERRORS:\n" + err.Error()
		} else {
			c <- "--------------\n" + "host: " + h + "\n" + string(response)
		}
		// workers <- c
	}()
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%s", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//func init() {
//
//	Init_log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
//}

func main() {
	var port int
	var dstHosts arrayFlags
	//var dstHostsString, command string
	var command string

	// initiate log levels
	InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	flag.IntVar(&port, "p", 22, "an int")
	flag.Var(&dstHosts, "h", `destination addresses
	Usage: -hosts=host:port
	Examples:
	-h=127.{1,2,3}               -> [ 127.1:22, 127.2:22, 127.3:22 ]
	-h=127.{1,2:2222,3}          -> [ 127.1:22, 127.2:2222, 127.3:22 ]
	-h=127.{1:22,2:22,4} -p 2222 -> [ 127.1:22, 127.2:22, 127.4:2222 ]`)
	flag.StringVar(&command, "c", "hostname", "bash command you want to execute on the destination hosts")
	flag.Parse()

	if len(dstHosts) == 0 {
		Error.Fatal("Destination Host is not defined")
	} else {
		Info.Println("dstHosts: ", dstHosts)
	}

	if command == "" {
		Error.Fatal("command can't be empty")
	} else {
		Info.Println("command: ", command)
	}
	Info.Println("default port: ", port)

	conf := makeUserConf()

	//// for test purposes (avoid env variables)
	//command = "sleep 3; hostname"

	var workers = make(chan chan string)

	for _, dstHost := range dstHosts {
		if !strings.Contains(dstHost, ":") {
			dstHost += ":" + strconv.Itoa(port)
		}
		//workers = append(workers, execCommand(dstHost, command, conf))
		execCommand(dstHost, command, conf, workers)
	}
	for i := 0; i < len(dstHosts); i++ {
		fmt.Println(<-(<-workers))
	}
}
