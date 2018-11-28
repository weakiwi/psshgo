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

psshgo default private key path is `/root/.ssh/id_rsa `, you can specific your own path with env `PKE`

## exec command on multiple hosts
host1 file content(this hosts list support below formats: user@10.102.10.2, user@10.102.10.2:2233, 
and default user is root,default port is 22)
```
10.102.10.2
10.102.10.3
10.102.10.4
```

and then

```bash
root@alpha: /home/weakiwi/code/psshgo/release master!
 # ./psshgo ssh -hf host1 -c "uname -r"                                                                                       [14:33:56]
2018/11/28 14:34:04 dial ssh success
2018/11/28 14:34:04 dial ssh success
2018/11/28 14:34:04 dial ssh success
2018/11/28 14:34:04 [10.102.10.3]
4.4.0-31-generic

2018/11/28 14:34:04 [10.102.10.2]
4.4.0-31-generic

2018/11/28 14:34:04 [10.102.10.4]
4.4.0-139-generic
```


## transfer file to multiple hosts
host1 file content
```
10.102.10.2
10.102.10.3
10.102.10.4
```

and then

```bash
root@alpha: /home/weakiwi/code/psshgo/release master!
 # ./psshgo scp -hf host1 -s host1 -d /home/host1                                                                             [14:36:32]
2018/11/28 14:36:40 dbb3f26dca4c95216dc24e833cb5f26a  host1
2018/11/28 14:36:42 dial ssh success
2018/11/28 14:36:42 dial ssh success
2018/11/28 14:36:42 dial ssh success
2018/11/28 14:36:43 [10.102.10.2]
dbb3f26dca4c95216dc24e833cb5f26a  /home/host1

2018/11/28 14:36:43 [10.102.10.3]
dbb3f26dca4c95216dc24e833cb5f26a  /home/host1

2018/11/28 14:36:43 [10.102.10.4]
dbb3f26dca4c95216dc24e833cb5f26a  /home/host1
```

## series of comannds and files

playbook example:

```ini
#host lists that you want to apply below tasks to
hosts=10.102.10.1,10.102.10.2,10.102.10.3,10.102.10.4

#task1 scp current folder playbook.ini to hosts /home/playbook.ini
[task1]
type=scp
src=playbook.ini
dst=/home/playbook.ini

#task2 exec 'cat /home/playbook.ini' on hosts
[task2]
type=ssh
command="cat /home/playbook.ini"
```

and then

```bash
root@alpha: /home/weakiwi/code/psshgo/release master!
 # ./psshgo pini -i playbook.ini                                                                                              [14:39:17]
2018/11/28 14:39:21 #######start  task1  ########
2018/11/28 14:39:21 51f9f8511a59d293589b260c44d26551  playbook.ini
2018/11/28 14:39:21
2018/11/28 14:39:24 dial ssh success
2018/11/28 14:39:24 dial ssh success
2018/11/28 14:39:24 dial ssh success
2018/11/28 14:39:25
2018/11/28 14:39:25 [10.102.10.1]

2018/11/28 14:39:25 [10.102.10.3]
51f9f8511a59d293589b260c44d26551  /home/playbook.ini

2018/11/28 14:39:25 #######start  task2  ########
2018/11/28 14:39:25 sshexecWithoutConnect sshclient.Cmd error:  ssh: handshake failed: EOF
2018/11/28 14:39:25
2018/11/28 14:39:25 [10.102.10.3]
#host lists that you want to apply below tasks to
hosts=10.102.10.1,10.102.10.2,10.102.10.3,10.102.10.4

#task1 scp current folder playbook.ini to hosts /home/playbook.ini
[task1]
type=scp
src=playbook.ini
dst=/home/playbook.ini

#task2 exec 'cat /home/playbook.ini' on hosts
dst=/home/playbook.ini

#task2 exec 'cat /home/playbook.ini' on hosts
[task2]
type=ssh
command="cat /home/playbook.ini"

2018/11/28 14:40:09 [10.102.10.2]
#host lists that you want to apply below tasks to
hosts=10.102.10.2,10.102.10.3,10.102.10.4

#task1 scp current folder playbook.ini to hosts /home/playbook.ini
[task1]
type=scp
src=playbook.ini
dst=/home/playbook.ini

#task2 exec 'cat /home/playbook.ini' on hosts
[task2]
type=ssh
command="cat /home/playbook.ini"

2018/11/28 14:40:09 [10.102.10.4]
#host lists that you want to apply below tasks to
hosts=10.102.10.2,10.102.10.3,10.102.10.4

#task1 scp current folder playbook.ini to hosts /home/playbook.ini
[task1]
type=scp
src=playbook.ini
dst=/home/playbook.ini

#task2 exec 'cat /home/playbook.ini' on hosts
[task2]
type=ssh
command="cat /home/playbook.ini"
```

## compile by myself
1. go get -u github.com/weakiwi/gosshtool
2. go get -u github.com/urfave/cli
3. ./build.sh

## TODO
1. 一些函数没做错误抛出，包括pssh和gosshtool中的。比如主机不可达时会直接空指针，这个是需要修复的
2. 希望能做成c/s模式，这样执行playbook也能更加灵活。可以增加诸如hostname-filter，quota-filter，执行成功验证等。
3. ini的格式不够灵活，希望能支持yaml吧
