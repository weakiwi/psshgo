package pssh

import (
        "github.com/urfave/cli"
    	"bufio"
    	"fmt"
    	"io"
    	"os"
    	"strings"
		"github.com/scottkiss/gosshtool"
        "io/ioutil"
        "log"
        "sync"
        "crypto/md5"
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
	sshCmd = cli.Command{
		Name:   "ssh",
		Usage:  "-hf hostfile -c command",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: pssh,
		Flags: []cli.Flag{
			hostfileFlag,
            commandfileFlag,
		},
	}
	scpCmd = cli.Command{
		Name:   "scp",
		Usage:  "-hf hostfile -s srcfile -d destfile",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: pscp,
		Flags: []cli.Flag{
			hostfileFlag,
            sourcefileFlag,
            destfileFlag,
		},
	}
	iniCmd = cli.Command{
		Name:   "ini",
		Usage:  "-hf hostfile -i inifile",
		//Usage:  "psshgo -hf hostfile  <\"cmds\" | cmdsfile>",
		Action: ini,
		Flags: []cli.Flag{
			hostfileFlag,
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
    fmt.Printf("%x  %s\n", h.Sum(nil), srcfile)
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
 	    var myconfig sshconfig
        if strings.Contains(string(line), "@") && strings.Contains(string(line), ":") {
                s := strings.Split(string(line), "@")
                myconfig.user = s[0]
                s1 := strings.Split(s[1], ":")
                myconfig.address = s1[0]
                myconfig.port = s1[1]
        } else if strings.Contains(string(line), ":") == false  && strings.Contains(string(line), "@"){
                s := strings.Split(string(line), "@")
                myconfig.user = s[0]
                myconfig.address = s[1]
                myconfig.port = "22"
        } else if strings.Contains(string(line), "@") == false && strings.Contains(string(line), ":"){
                myconfig.user = "root"
                s := strings.Split(string(line), ":")
                myconfig.address = s[0]
                myconfig.port = s[1]
        } else {
                myconfig.user = "root"
                myconfig.address = strings.Replace(string(line), "\n", "", -1)
                myconfig.port = "22"
        }
        result_sshconfig = append(result_sshconfig, myconfig)
   }
   return result_sshconfig, nil
}
func ini(c *cli.Context) {
    fmt.Println("hello this ini")
}
func pscp(c *cli.Context) {
    hostfile := mustGetStringVar(c, "hf")
    srcfile := mustGetStringVar(c, "s")
    destfile := mustGetStringVar(c, "d")
    //var t *testing.T
    fi, err := os.Open(hostfile)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    counter := ComputeLine(hostfile)
    done := make(chan string, counter)
    myconfigs, err := parseHostfile(hostfile)
    if err != nil {
        log.Fatalf("pscp.parseHostfile err: %v", err)
    }
    for i := range myconfigs{
        waitgroup.Add(1)
        go scpexec(&myconfigs[i], srcfile, destfile, done)
    }
	md5File(srcfile)
    waitgroup.Wait()
	for v := range done {
	    fmt.Println(v)
	    if len(done) <= 0 { // 如果现有数据量为0，跳出循环
            close(done)
	    }
	}
    return
}
func ComputeLine(path string)(num int){
    f,err := os.Open(path)
    if nil != err{
        log.Println(err)
        return
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for{
        _,err := r.ReadString('\n')
        if io.EOF == err || nil != err{
            break
        }
        num += 1
    }
    return
}
func pssh(c *cli.Context) {
    hostfile := mustGetStringVar(c, "hf")
    command := mustGetStringVar(c, "c")

    //var t *testing.T
    fi, err := os.Open(hostfile)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    br := bufio.NewReader(fi)
    counter := ComputeLine(hostfile)
    done := make(chan string, counter)
    for {
        line, err := br.ReadString('\n')
	    var myconfig sshconfig
        if err != nil || err == io.EOF {
            break
        }
        if strings.Contains(string(line), "@") && strings.Contains(string(line), ":") {
                s := strings.Split(string(line), "@")
                myconfig.user = s[0]
                s1 := strings.Split(s[1], ":")
                myconfig.address = s1[0]
                myconfig.port = s1[1]
        } else if strings.Contains(string(line), ":") == false  && strings.Contains(string(line), "@"){
                s := strings.Split(string(line), "@")
                myconfig.user = s[0]
                myconfig.address = s[1]
                myconfig.port = "22"
        } else if strings.Contains(string(line), "@") == false && strings.Contains(string(line), ":"){
                myconfig.user = "root"
                s := strings.Split(string(line), ":")
                myconfig.address = s[0]
                myconfig.port = s[1]
        } else {
                myconfig.user = "root"
                myconfig.address = strings.Replace(string(line), "\n", "", -1)
                myconfig.port = "22"
        }
        waitgroup.Add(1)
        go sshexec(&myconfig, command, done)
    }
    waitgroup.Wait()
	for v := range done {
	    fmt.Println(v)
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
		User:     sc.user,
		Privatekey: pkey,
        Host:     sc.address,
	}
    sshclient := gosshtool.NewSSHClient(config2)
    stdout, stderr, _, err := sshclient.Cmd(command, nil, nil, 0)
	if err != nil {
        waitgroup.Done()
		done <- fmt.Sprintf(stderr)
        return
	}
    waitgroup.Done()
	done <- fmt.Sprintf("%s[%s]%s\n%s", CLR_R, sc.address, CLR_N,stdout)
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
		User:     sc.user,
		Privatekey: pkey,
        Host:     sc.address,
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
            fmt.Printf(stderr)
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
type sshconfig struct {
	user string
	address string
    port string
}

func mustGetStringVar(c *cli.Context, key string) string {
	v := strings.TrimSpace(c.String(key))
	if v == "" {
		errExit(1, "%s must be provided", key)
	}
	return v
}

func errExit(code int, format string, val ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", val...)
	os.Exit(code)
}

