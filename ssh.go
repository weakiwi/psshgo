package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"golang.org/x/crypto/ssh"
)

func cliSSH(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Println("Enter a quote enclosed command or file containing commands")
		os.Exit(1)
	}

	// get the username, hosts, and commands from cli
	userName := c.String("user")

	hosts := loadHosts(c.String("host"))

	command := c.Args()[0]
	// if its a file, parse
	if _, err := os.Stat(command); err == nil {
		command = commandFromFile(command)
	}

	privateKey := decryptKeyFile(path.Join(home(), ".ssh/id_rsa"))
	clientConfig := makeClientConfig(userName, privateKey)

	done := make(chan Result)
	liftoff := time.Now()
	for _, host := range hosts {
		go runCommandOnHost(clientConfig, command, host, done)
	}

	dir, file := timeToDirFile(liftoff)
	for _, host := range hosts {
		res := <-done
		spl := strings.Split(host, ":")
		addr := spl[0]
		var result []byte
		if res.err != nil {
			result = []byte(res.err.Error())
		} else {
			result = res.buf.Bytes()
		}
		switch c.String("out") {
		case "stdout":
			fmt.Println(string(result))
		default:
			err := ioutil.WriteFile(path.Join(dir, file+"_"+addr), result, 0600)
			if err != nil {
				fmt.Println("Error writing host to file: ", host, err)
			}

		}
	}
}

type Result struct {
	buf *bytes.Buffer
	err error
}

func runCommandOnHost(config *ssh.ClientConfig, cmd, host string, done chan Result) {
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		done <- Result{err: fmt.Errorf("Failed to dial: %s %s", host, err.Error())}
		return
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		done <- Result{err: fmt.Errorf("Failed to create session: %s", err.Error())}
		return
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	b := new(bytes.Buffer)
	session.Stdout = b
	if err := session.Run(cmd); err != nil {
		done <- Result{err: fmt.Errorf("Failed to run: %s", err.Error())}
	}
	done <- Result{buf: b}
}

func commandFromFile(command string) string {
	b, err := ioutil.ReadFile(command)
	ifExit(err)

	buf := new(bytes.Buffer)
	spl := strings.Split(string(b), "\n")
	for _, s := range spl {
		fmt.Fprintf(buf, s+"; ")
	}

	return string(buf.Bytes())
}
