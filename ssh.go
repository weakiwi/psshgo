package main

import (
	"os"
	"io/ioutil"
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