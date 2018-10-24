package dbx_sharing

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/dbx_api"
	"github.com/watermint/toolbox/dbx_api/dbx_rpc"
	"time"
)

type SharedLink struct {
	SharedLinkId  string          `json:"shared_link_id"`
	AsMemberId    string          `json:"as_member_id"`
	AsMemberEmail string          `json:"as_member_email"`
	Url           string          `json:"url"`
	Link          json.RawMessage `json:"link"`
}

func (a *SharedLink) UpdateExpire(c *dbx_api.Context, newExpire time.Time) (newLink *SharedLink, annotation dbx_api.ErrorAnnotation, err error) {
	url := gjson.Get(string(a.Link), "url").String()

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
			Expires: dbx_api.RebaseTimeForAPI(newExpire).Format(dbx_api.DateTimeFormat),
		},
	}

	req := dbx_rpc.RpcRequest{
		Endpoint:   "sharing/modify_shared_link_settings",
		Param:      up,
		AsMemberId: a.AsMemberId,
	}
	res, ea, err := req.Call(c)
	if ea.IsFailure() {
		return nil, ea, err
	}

	newLink, ea, err = ParseSharedLink(gjson.Get(res.Body, ""))
	if ea.IsFailure() {
		return nil, ea, err
	}
	newLink.AsMemberId = a.AsMemberId
	newLink.AsMemberEmail = a.AsMemberEmail

	return newLink, dbx_api.Success, nil
}

func ParseSharedLink(res gjson.Result) (link *SharedLink, annotation dbx_api.ErrorAnnotation, err error) {
	linkId := res.Get("id")
	if !linkId.Exists() {
		err = errors.New("required field `id` not found")
		annotation = dbx_api.ErrorAnnotation{
			ErrorType: dbx_api.ErrorUnexpectedDataType,
			Error:     err,
		}
		return nil, annotation, err
	}
	url := res.Get("url")
	if !linkId.Exists() {
		err = errors.New("required field `url` not found")
		annotation = dbx_api.ErrorAnnotation{
			ErrorType: dbx_api.ErrorUnexpectedDataType,
			Error:     err,
		}
		return nil, annotation, err
	}

	s := &SharedLink{
		SharedLinkId: linkId.String(),
		Url:          url.String(),
		Link:         json.RawMessage(res.Raw),
	}
	return s, dbx_api.Success, nil
}

type SharedLinkList struct {
	Path          string
	AsMemberId    string
	AsMemberEmail string
	OnError       func(annotation dbx_api.ErrorAnnotation) bool
	OnEntry       func(link *SharedLink) bool
}

func (a *SharedLinkList) List(c *dbx_api.Context) bool {
	type ListParam struct {
		Path string `json:"path,omitempty"`
	}
	lp := ListParam{
		Path: a.Path,
	}
	list := dbx_rpc.RpcList{
		EndpointList:         "sharing/list_shared_links",
		EndpointListContinue: "sharing/list_shared_links",
		AsMemberId:           a.AsMemberId,
		UseHasMore:           true,
		ResultTag:            "links",
		OnError:              a.OnError,
		OnEntry: func(link gjson.Result) bool {
			if a.OnEntry == nil {
				return true
			}

			s, ea, _ := ParseSharedLink(link)
			if ea.IsSuccess() {
				s.AsMemberId = a.AsMemberId
				s.AsMemberEmail = a.AsMemberEmail
				return a.OnEntry(s)
			}

			if a.OnError != nil {
				return a.OnError(ea)
			}
			return false
		},
	}

	return list.List(c, lp)
}