package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetsubinfo"
)

const (
	ProtoIDQotGetSubInfo = 3003 //Qot_GetSubInfo	获取订阅信息
)

// 获取订阅信息
func (api *FutuAPI) QuerySubscription(ctx context.Context, isAll bool) (*Subscription, error) {
	ch := make(qotgetsubinfo.ResponseChan)
	if err := api.get(ProtoIDQotGetSubInfo, &qotgetsubinfo.Request{C2S: &qotgetsubinfo.C2S{
		IsReqAllConn: &isAll,
	}}, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return subscriptionFromPB(resp.GetS2C()), result(resp)
	}
}

type Subscription struct {
	ConnSubInfos   []*ConnSubInfo
	TotalUsedQuota int32
	RemainQuota    int32
}

func subscriptionFromPB(pb *qotgetsubinfo.S2C) *Subscription {
	if pb == nil {
		return nil
	}
	sub := Subscription{
		TotalUsedQuota: pb.GetTotalUsedQuota(),
		RemainQuota:    pb.GetRemainQuota(),
	}
	if list := pb.GetConnSubInfoList(); list != nil {
		sub.ConnSubInfos = make([]*ConnSubInfo, len(list))
		for i, v := range list {
			sub.ConnSubInfos[i] = connSubInfoFromPB(v)
		}
	}
	return &sub
}

type ConnSubInfo struct {
	SubInfos  []*SubInfo
	UsedQuota int32
	IsOwnData bool
}

func connSubInfoFromPB(pb *qotcommon.ConnSubInfo) *ConnSubInfo {
	if pb == nil {
		return nil
	}
	info := ConnSubInfo{
		UsedQuota: pb.GetUsedQuota(),
		IsOwnData: pb.GetIsOwnConnData(),
	}
	if list := pb.GetSubInfoList(); list != nil {
		info.SubInfos = make([]*SubInfo, len(list))
		for i, v := range list {
			info.SubInfos[i] = subInfoFromPB(v)
		}
	}
	return &info
}

type SubInfo struct {
	SubType    qotcommon.SubType
	Securities []*Security
}

func subInfoFromPB(pb *qotcommon.SubInfo) *SubInfo {
	if pb == nil {
		return nil
	}
	info := SubInfo{
		SubType: qotcommon.SubType(pb.GetSubType()),
	}
	if list := pb.GetSecurityList(); list != nil {
		info.Securities = make([]*Security, len(list))
		for i, v := range list {
			info.Securities[i] = securityFromPB(v)
		}
	}
	return &info
}
