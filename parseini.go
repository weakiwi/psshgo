package main

import (
        "os"
        "log"
        "github.com/go-ini/ini"
)

func main() {
        cfg, err := ini.Load("playbook.ini")
        if err != nil {
            log.Fatalf("open ini err: %v", err)
            os.Exit(1)
        }
        secs := cfg.Sections()
        for i := range secs {
                ini_type := secs[i].Key("type").String()
                if ini_type == "ssh" {
                    log.Println("command is : %v", secs[i].Key("command"))
                } else if ini_type == "scp" {
                    log.Println("dst is : %v", secs[i].Key("dst"))
                    log.Println("src is : %v", secs[i].Key("src"))
                }
        }
}
