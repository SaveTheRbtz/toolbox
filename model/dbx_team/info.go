package dbx_team

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/model/dbx_api"
	"github.com/watermint/toolbox/model/dbx_rpc"
)

type TeamInfo struct {
	TeamId string          `json:"team_id"`
	Info   json.RawMessage `json:"info"`
}

type TeamInfoList struct {
	OnError func(err error) bool
	OnEntry func(info *TeamInfo) bool
}

func (t *TeamInfoList) List(c *dbx_api.Context) bool {
	req := dbx_rpc.RpcRequest{
		Endpoint: "team/get_info",
	}
	res, err := req.Call(c)
	if err != nil {
		return t.OnError(err)
	}

	teamId := gjson.Get(res.Body, "team_id").String()
	team := &TeamInfo{
		TeamId: teamId,
		Info:   json.RawMessage(res.Body),
	}

	if t.OnEntry != nil {
		return t.OnEntry(team)
	}
	return true
}
