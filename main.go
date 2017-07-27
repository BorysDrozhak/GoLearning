package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	user := "bdrozhak"
	home := "/Users/" + user
	private_key_path := "/.ssh/id_rsa"
	ssh_password := "ypur password"
	dstHost := "127.0.0.1:22"
	command := "env"

	key, err := ioutil.ReadFile(home + private_key_path)
	if err != nil {
		log.Fatalf("unable to read private key: %v", key)
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Println("Private key is encrypted")
		passPhrase := os.Getenv("PASSPHRASE")
		if passPhrase == "" {
			log.Println("PASSPHRASE is not set as env variable")
			// TODO: ask passphrase in case
		}
		rawKey, err := ssh.ParseRawPrivateKeyWithPassphrase(key, []byte(passPhrase))
		signer, err = ssh.NewSignerFromKey(rawKey)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(ssh_password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", dstHost, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	fmt.Println(b.String())
}
