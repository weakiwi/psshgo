# psshgo

golang paralle ssh

```
NAME:
   psshgo - psshgo

USAGE:
   psshgo [global options] command [command options] [arguments...]

VERSION:
   1.0

AUTHOR:
   weakiwi <dengyi0215@gmail.com>

COMMANDS:
     ssh      -hf hostfile -c command
     scp      -hf hostfile -s srcfile -d destfile
     pini     -i inifile
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

playbook example:
```
[group1]
type=hosts
ips=10.102.10.1,10.102.10.2,10.102.10.3,10.102.10.4
[task1]
type=scp
src=playbook.ini
dst=/home/playbook.ini
[task2]
type=ssh
command="cat /home/playbook.ini"
```
