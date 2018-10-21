package sharedlink

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/dbx_api"
	"github.com/watermint/toolbox/dbx_api/dbx_rpc"
	"github.com/watermint/toolbox/workflow"
	"time"
)

type WorkerSharedLinkUpdateExpires struct {
	workflow.SimpleWorkerImpl
	Api      *dbx_api.Context
	NextTask string
	Days     int
}

type ContextSharedLinkUpdateExpiresResult struct {
	AsMemberId    string          `json:"as_member_id"`
	AsMemberEmail string          `json:"as_member_email"`
	SharedLinkId  string          `json:"shared_link_id"`
	ExpiresOld    string          `json:"expires_old"`
	ExpiresNew    string          `json:"expires_new"`
	Link          json.RawMessage `json:"link"`
}

func (w *WorkerSharedLinkUpdateExpires) Prefix() string {
	return "/sharedlink/update/expires"
}

func (w *WorkerSharedLinkUpdateExpires) Exec(task *workflow.Task) {
	tc := &ContextSharedLinkResult{}
	workflow.UnmarshalContext(task, tc)

	link := string(tc.Link)
	expires := gjson.Get(link, "expires").String()
	var origTime time.Time

	if expires != "" {
		var err error
		origTime, err = time.Parse(dbx_api.DateTimeFormat, expires)
		if err != nil {
			seelog.Warnf("SharedLinkId[%s] Unable to parse time [%s]", tc.SharedLinkId, expires)
			return
		}
	}

	targetExpire := dbx_api.RebaseTimeForAPI(time.Now().Add(time.Duration(w.Days*24) * time.Hour))
	seelog.Debugf("LinkId[%s] Link's expire time[%s] Target[%s]", tc.SharedLinkId, origTime.String(), targetExpire.String())
	if origTime.IsZero() || origTime.After(targetExpire) {
		w.update(targetExpire, origTime, tc, task)
	} else {
		seelog.Debugf("Skip LinkId[%s] Expire[%s]", tc.SharedLinkId, origTime.String())
	}
}

func (w *WorkerSharedLinkUpdateExpires) update(targetTime time.Time, origTime time.Time, tc *ContextSharedLinkResult, task *workflow.Task) {
	link := string(tc.Link)
	url := gjson.Get(link, "url").String()

	type SettingsParam struct {
		Expires string `json:"expires"`
	}
	type UpdateParam struct {
		Url      string        `json:"url"`
		Settings SettingsParam `json:"settings"`
	}

	up := UpdateParam{
		Url: url,
		Settings: SettingsParam{
			Expires: targetTime.Format(dbx_api.DateTimeFormat),
		},
	}

	oldTime := origTime.Format(dbx_api.DateTimeFormat)
	newTime := targetTime.Format(dbx_api.DateTimeFormat)

	seelog.Infof("Updating: SharedLinkId[%s] MemberEmail[%s]: Old[%s] -> New[%s]", tc.SharedLinkId, tc.AsMemberEmail, oldTime, newTime)

	req := dbx_rpc.RpcRequest{
		Endpoint:   "sharing/modify_shared_link_settings",
		Param:      up,
		AsMemberId: tc.AsMemberId,
	}
	res, ea, _ := req.Call(w.Api)
	if ea.IsFailure() {
		w.Pipeline.HandleGeneralFailure(ea)
		return
	}

	w.Pipeline.Enqueue(
		workflow.MarshalTask(
			w.NextTask,
			tc.SharedLinkId,
			ContextSharedLinkUpdateExpiresResult{
				AsMemberId:    tc.AsMemberId,
				AsMemberEmail: tc.AsMemberEmail,
				SharedLinkId:  tc.SharedLinkId,
				ExpiresOld:    oldTime,
				ExpiresNew:    newTime,
				Link:          json.RawMessage(res.Body),
			},
		),
	)
}
