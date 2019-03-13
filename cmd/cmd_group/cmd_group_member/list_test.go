package cmd_group_member

import (
	"github.com/watermint/toolbox/cmd"
	"io/ioutil"
	"os"
	"testing"
)

func TestCmdGroupMemberList_Exec(t *testing.T) {
	csvFile, err := ioutil.TempFile("", "group_member_list")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(csvFile.Name())

	cmd.CmdTest(t, NewCmdGroupMember(), []string{"list", "-report-path", csvFile.Name()})
}