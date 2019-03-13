package cmd_namespace_file

import (
	"github.com/watermint/toolbox/cmd"
	"io/ioutil"
	"os"
	"testing"
)

func TestCmdTeamNamespaceFileList_Exec(t *testing.T) {
	csvFile, err := ioutil.TempFile("", "member_list")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(csvFile.Name())

	cmd.CmdTest(t, NewCmdTeamNamespaceFile(), []string{"list", "-report-path", csvFile.Name()})
}