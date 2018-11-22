package pssh

import (
        "os"
        "log"
        "github.com/go-ini/ini"
)

func parseini(path string) (scp_playbooks []playbook_scp, ssh_playbooks []playbook_ssh, err error){
        cfg, err := ini.Load(path)
        if err != nil {
            log.Fatalf("open ini err: %v", err)
            return nil, nil, err
            os.Exit(1)
        }
        secs := cfg.Sections()
        for i := range secs {
                ini_type := secs[i].Key("type").String()
                if ini_type == "ssh" {
                    var my_playbook_ssh playbook_ssh
                    log.Println("command is : %v", secs[i].Key("command"))
                    my_playbook_ssh.command = secs[i].Key("command").String()
                    ssh_playbooks = append(ssh_playbooks, my_playbook_ssh)
                } else if ini_type == "scp" {
                    var my_playbook_scp playbook_scp
                    log.Println("dst is : %v", secs[i].Key("dst"))
                    log.Println("src is : %v", secs[i].Key("src"))
                    my_playbook_scp.dst = secs[i].Key("dst").String()
                    my_playbook_scp.src = secs[i].Key("src").String()
                    scp_playbooks = append(scp_playbooks, my_playbook_scp)
                }
        }
        return scp_playbooks, ssh_playbooks, nil
}

type playbook_ssh struct {
    command string,
}

type playbook_scp struct {
    dst string,
    src string,
}