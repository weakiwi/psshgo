package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"log"
	"github.com/weakiwi/gosshtool"
)
func make_a_connection(sc *sshconfig) (sshclient *gosshtool.SSHClient, err error) {
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
		User:       sc.user,
		Privatekey: pkey,
		Host:       sc.address,
	}
	sshclient = gosshtool.NewSSHClient(config2)
	return sshclient, nil
}

func sshexec_without_connect(sshclient *gosshtool.SSHClient, command string, done chan string) {
	stdout, stderr, _, err := sshclient.Cmd(command, nil, nil, 0)
	if err != nil {
		waitgroup.Done()
		log.Println("sshexec error is : ", err)
		done <- fmt.Sprintf(stderr)
		return
	}
	waitgroup.Done()
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, sshclient.SSHClientConfig.Host, CLR_N, stdout)
	return
}

func scpexec_without_connection(client *gosshtool.SSHClient, srcfile string, destfile string, done chan string) {
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
		waitgroup.Done()
		done <- fmt.Sprintf(stderr)
	}
	waitgroup.Done()
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, client.SSHClientConfig.Host, CLR_N, stdout)
	return
}