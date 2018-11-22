package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"log"
	"os"
	"strings"
)

func parseini(path string) (playbooks []playbook, err error) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("open ini err: %v", err)
		return nil, err
		os.Exit(1)
	}
	var my_playbook playbook
	var sshconfigs []sshconfig
	secs := cfg.Sections()
	for i := range secs {
        ini_type := secs[i].Key("type").String()
        my_playbook.name = secs[i].Name()
		if ini_type == "hosts" {
			ips := strings.Split(secs[i].Key("ips").String(), ",")
			for j := range ips {
				k, err := stringToSshconfig(ips[j])
				if err != nil {
					log.Fatalf("parseini error: %v", err)
				}
				sshconfigs = append(sshconfigs, k)
			}
			continue
		}
		// if operate type is ssh
		if secs[i].HasKey("command") == true {
			my_playbook.command = secs[i].Key("command").String()
			my_playbook.playbook_type = "ssh"
			my_playbook.dst = ""
			my_playbook.src = ""
			playbooks = append(playbooks, my_playbook)
			continue
			// if operate type is scp
		} else if secs[i].HasKey("dst") == true && secs[i].HasKey("src") {
			my_playbook.dst = secs[i].Key("dst").String()
			my_playbook.src = secs[i].Key("src").String()
			my_playbook.playbook_type = "scp"
			my_playbook.command = ""
			playbooks = append(playbooks, my_playbook)
			continue
		}
	}
	for i := range playbooks {
		playbooks[i].servers = sshconfigs
	}
	return playbooks, nil
}

func stringToSshconfig(line string) (myconfig sshconfig, err error) {
	if strings.Contains(string(line), "@") && strings.Contains(string(line), ":") {
		s := strings.Split(string(line), "@")
		myconfig.user = s[0]
		s1 := strings.Split(s[1], ":")
		myconfig.address = s1[0]
		myconfig.port = s1[1]
	} else if strings.Contains(string(line), ":") == false && strings.Contains(string(line), "@") {
		s := strings.Split(string(line), "@")
		myconfig.user = s[0]
		myconfig.address = s[1]
		myconfig.port = "22"
	} else if strings.Contains(string(line), "@") == false && strings.Contains(string(line), ":") {
		myconfig.user = "root"
		s := strings.Split(string(line), ":")
		myconfig.address = s[0]
		myconfig.port = s[1]
	} else {
		myconfig.user = "root"
		myconfig.address = strings.Replace(string(line), "\n", "", -1)
		if myconfig.address == "" {
			log.Fatalf("stringToSshconfig error: line is blank!")
			return myconfig, fmt.Errorf("stringToSshconfig error: line is blank!")
		}
		myconfig.port = "22"
	}
	return myconfig, nil

}
func parseHostfile(hostfile string) (result_sshconfig []sshconfig, err error) {
	fi, err := os.Open(hostfile)
	if err != nil {
		fmt.Printf("parseHostfile.Open Error: %s\n", err)
		return nil, err
	}
	br := bufio.NewReader(fi)
	for {
		line, err := br.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		myconfig, err := stringToSshconfig(line)
		if err != nil {
			log.Fatalf("parseHostfile.stringToSshconfig error: %v", err)
			return nil, err
		}
		result_sshconfig = append(result_sshconfig, myconfig)
	}
	return result_sshconfig, nil
}

type sshconfig struct {
	user    string
	address string
	port    string
}

type playbook struct {
    name          string
	playbook_type string
	src           string
	dst           string
	command       string
	servers       []sshconfig
}
