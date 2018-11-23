package test

import (
	"testing"
)

func Test_Parseini(t *testing.T) {
	playbooks, err := parseini("playbook.ini")
	if err != nil {
		t.Error("parseini error: %v", err)
	}
	my_playbook_1 := playbooks[0]
	my_playbook_2 := playbooks[1]
	if my_playbook_1.type == "scp" and my_playbook_1.type == "ssh" {
		t.Log("parseini get type pass")
	}
	if my_playbook_1.src == "playbook.ini" && my_playbook_1.dst == "/home/playbook.ini" && my_playbook_2.command === "cat /home/playbook.ini" {
		t.Log("get other paraments pass")
	}
}