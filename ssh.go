package main

import (
	"fmt"
	"github.com/weakiwi/gosshtool"
	"io/ioutil"
	"log"
	"os"
)

func makeConnection(sc *sshconfig) (sshclient *gosshtool.SSHClient, err error) {
	pkey := os.Getenv("PKEY")
	if pkey == "" {
		pkey = "/root/.ssh/id_rsa"
	}
	key, err := ioutil.ReadFile(pkey)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
		return nil, err
	}
	pkey = string(key)
	config2 := &gosshtool.SSHClientConfig{
		User:              sc.user,
		Privatekey:        pkey,
		Host:              sc.address,
		DialTimeoutSecond: 5,
	}
	sshclient = gosshtool.NewSSHClient(config2)
	return sshclient, nil
}

func sshexecWithoutConnect(sshclient *gosshtool.SSHClient, command string, done chan string) {
	defer waitgroup.Done()
	stdout, stderr, _, err := sshclient.Cmd(command, nil, nil, 0)
	if err != nil {
		log.Println("sshexecWithoutConnect sshclient.Cmd error: ", err)
		done <- fmt.Sprintf(stderr)
		return
	}
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, sshclient.SSHClientConfig.Host, CLR_N, stdout)
	return
}

func scpexecWithoutConnection(client *gosshtool.SSHClient, srcfile string, destfile string, done chan string) {
	defer waitgroup.Done()
	f, err := os.Open(srcfile)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	stdout, stderr, err := client.TransferData(destfile, data)
	if err != nil {
		log.Printf(stderr)
	}
	stdout, stderr, _, err = client.Cmd("md5sum "+destfile, nil, nil, 0)
	if err != nil {
		done <- fmt.Sprintf(stderr)
	}
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, client.SSHClientConfig.Host, CLR_N, stdout)
	return
}
