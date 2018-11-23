package main

import (
	"bufio"
	"crypto/md5"
	"github.com/scottkiss/gosshtool"
	"github.com/urfave/cli"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	pssh_logfile := os.Getenv("P_LOGFILE")
	if pssh_logfile != "" {
		f, _ := os.OpenFile(string(pssh_logfile), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND,0755)
		fmt.Println("main pssh_logfile: ", f)
		os.Stdout = f
		os.Stderr = f
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

func md5File(srcfile string) {
	file, err := os.Open(srcfile)
	if err != nil {
		panic(err)
	}

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return
	}
	log.Printf("%x  %s\n", h.Sum(nil), srcfile)
}

func pini(c *cli.Context) {
	inifile := mustGetStringVar(c, "i")
	playbooks, err := parseini(inifile)
	if err != nil {
		log.Fatalf("pini error: %v", err)
		os.Exit(1)
	}
	for i := range playbooks {
		log.Println("#######start ", playbooks[i].name, " ########")
		if playbooks[i].playbook_type == "scp" {
			pscpexec(playbooks[i].servers, playbooks[i].src, playbooks[i].dst)
		} else if playbooks[i].playbook_type == "ssh" {
			psshexec(playbooks[i].servers, playbooks[i].command)
		}
	}
}
func pscpexec(servers []sshconfig, srcfile string, destfile string) {
	counter := len(servers)
	done := make(chan string, counter)
	for i := range servers {
		waitgroup.Add(1)
		go scpexec(&servers[i], srcfile, destfile, done)
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
func ComputeLine(path string) (num int) {
	f, err := os.Open(path)
	if nil != err {
		log.Println(err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		_, err := r.ReadString('\n')
		if io.EOF == err || nil != err {
			break
		}
		num += 1
	}
	return
}

func psshexec(servers []sshconfig, command string) {
	counter := len(servers)
	done := make(chan string, counter)
	for i := range servers {
		waitgroup.Add(1)
		go sshexec(&servers[i], command, done)
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
	log.Println("pssh myconfigs are : %v", myconfigs)
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
	pkey := os.Getenv("PKEY")
	if pkey == "" {
		pkey = "/root/.ssh/id_rsa"
	}
	key, err := ioutil.ReadFile(pkey)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}
	pkey = string(key)
	config2 := &gosshtool.SSHClientConfig{
		User:       sc.user,
		Privatekey: pkey,
		Host:       sc.address,
	}
	sshclient := gosshtool.NewSSHClient(config2)
	stdout, stderr, _, err := sshclient.Cmd(command, nil, nil, 0)
	if err != nil {
		waitgroup.Done()
		log.Println("sshexec error is : %v", err)
		done <- fmt.Sprintf(stderr)
		return
	}
	waitgroup.Done()
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, sc.address, CLR_N, stdout)
	return

}

func scpexec(sc *sshconfig, srcfile string, destfile string, done chan string) {
	pkey := os.Getenv("PKEY")
	if pkey == "" {
		pkey = "/root/.ssh/id_rsa"
	}
	key, err := ioutil.ReadFile(pkey)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}
	pkey = string(key)
	config2 := &gosshtool.SSHClientConfig{
		User:       sc.user,
		Privatekey: pkey,
		Host:       sc.address,
	}
	client := gosshtool.NewSSHClient(config2)
	f, err := os.Open(srcfile)
	if err != nil {
		return
	}
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
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, sc.address, CLR_N, stdout)
	return
}

func mustGetStringVar(c *cli.Context, key string) string {
	v := strings.TrimSpace(c.String(key))
	if v == "" {
		log.Fatalf("%s must be provided", key)
	}
	return v
}