package main

import (
	"github.com/urfave/cli"
	"github.com/weakiwi/gosshtool"
	"log"
	"os"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Name = "psshgo"
	app.Usage = "psshgo"
	app.Version = "1.0"
	app.Author = "weakiwi"
	app.Email = "dengyi0215@gmail.com"
	app.Commands = []cli.Command{
		sshCmd,
		scpCmd,
		iniCmd,
	}

	app.Run(os.Args)
}

const CLR_R = "\x1b[31;1m"
const CLR_N = "\x1b[0m"

var (
	waitgroup sync.WaitGroup
	sshCmd    = cli.Command{
		Name:  "ssh",
		Usage: "-hf hostfile -c command",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: pssh,
		Flags: []cli.Flag{
			hostfileFlag,
			commandfileFlag,
		},
	}
	scpCmd = cli.Command{
		Name:  "scp",
		Usage: "-hf hostfile -s srcfile -d destfile",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: pscp,
		Flags: []cli.Flag{
			hostfileFlag,
			sourcefileFlag,
			destfileFlag,
		},
	}
	iniCmd = cli.Command{
		Name:  "pini",
		Usage: "-i inifile",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: pini,
		Flags: []cli.Flag{
			inifileFlag,
		},
	}

	hostfileFlag = cli.StringFlag{
		Name:  "hf",
		Usage: "host or file containing host names, one per line",
	}
	commandfileFlag = cli.StringFlag{
		Name:  "c",
		Usage: "command",
	}
	sourcefileFlag = cli.StringFlag{
		Name:  "s",
		Usage: "source file",
	}
	destfileFlag = cli.StringFlag{
		Name:  "d",
		Usage: "destination file",
	}
	inifileFlag = cli.StringFlag{
		Name:  "i",
		Usage: "ini file",
	}
)

func pini(c *cli.Context) {
	inifile := mustGetStringVar(c, "i")
	playbooks, err := parseini(inifile)
	if err != nil {
		log.Fatalf("pini error: %v", err)
		os.Exit(1)
	}
	var sc_group []*gosshtool.SSHClient
	if playbooks[0].servers == nil {
		log.Fatalf("playbooks format error")
		os.Exit(1)
	}
	for j := range playbooks[0].servers {
		tmp_conn, err := make_a_connection(&playbooks[0].servers[j])
		if err != nil {
			log.Fatalf("make_a_connection error: %v", err)
			os.Exit(1)
		}
		sc_group = append(sc_group, tmp_conn)
	}
	for i := range playbooks {
		log.Println("#######start ", playbooks[i].name, " ########")
		if playbooks[i].playbook_type == "scp" {
			pscpexec(sc_group, playbooks[i].src, playbooks[i].dst)
		} else if playbooks[i].playbook_type == "ssh" {
			psshexec(sc_group, playbooks[i].command)
		}
	}
}
func pscpexec(servers []*gosshtool.SSHClient, srcfile string, destfile string) {
	counter := len(servers)
	done := make(chan string, counter)
	for i := range servers {
		waitgroup.Add(1)
		go scpexec_without_connection(servers[i], srcfile, destfile, done)
	}
	md5File(srcfile)
	waitgroup.Wait()
	for v := range done {
		log.Println(v)
		if len(done) <= 0 {
			close(done)
		}
	}
}
func pscp(c *cli.Context) {
	hostfile := mustGetStringVar(c, "hf")
	srcfile := mustGetStringVar(c, "s")
	destfile := mustGetStringVar(c, "d")
	//var t *testing.T
	counter := ComputeLine(hostfile)
	done := make(chan string, counter)
	myconfigs, err := parseHostfile(hostfile)
	if err != nil {
		log.Fatalf("pscp.parseHostfile err: %v", err)
	}
	for i := range myconfigs {
		waitgroup.Add(1)
		go scpexec(&myconfigs[i], srcfile, destfile, done)
	}
	md5File(srcfile)
	waitgroup.Wait()
	for v := range done {
		log.Println(v)
		if len(done) <= 0 { // 如果现有数据量为0，跳出循环
			close(done)
		}
	}
	return
}

func psshexec(servers []*gosshtool.SSHClient, command string) {
	counter := len(servers)
	done := make(chan string, counter)
	for i := range servers {
		waitgroup.Add(1)
		go sshexec_without_connect(servers[i], command, done)
	}
	waitgroup.Wait()
	for v := range done {
		log.Println(v)
		if len(done) <= 0 {
			close(done)
		}
	}
}
func pssh(c *cli.Context) {
	hostfile := mustGetStringVar(c, "hf")
	command := mustGetStringVar(c, "c")

	//var t *testing.T
	counter := ComputeLine(hostfile)
	done := make(chan string, counter)
	myconfigs, err := parseHostfile(hostfile)
	if err != nil {
		log.Fatalf("sshexec.parseHostfile err: %v", err)
		os.Exit(1)
	}
	for i := range myconfigs {
		waitgroup.Add(1)
		go sshexec(&myconfigs[i], command, done)
	}
	waitgroup.Wait()
	for v := range done {
		log.Println(v)
		if len(done) <= 0 { // 如果现有数据量为0，跳出循环
			close(done)
		}
	}
}

func sshexec(sc *sshconfig, command string, done chan string) {
	sshclient, err := make_a_connection(sc)
	if err != nil {
		log.Fatalf("sshexec.make_a_connection error: %v", err)
	}
	sshexec_without_connect(sshclient, command, done)
}

func scpexec(sc *sshconfig, srcfile string, destfile string, done chan string) {
	sshclient, err := make_a_connection(sc)
	if err != nil {
		log.Fatalf("sshexec.make_a_connection error: %v", err)
	}
	scpexec_without_connection(sshclient, srcfile, destfile, done)
}
