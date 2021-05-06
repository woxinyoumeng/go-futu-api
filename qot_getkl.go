package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetkl"
)

const (
	ProtoIDQotGetKL = 3006 //Qot_GetKL	获取 K 线
)

// 获取实时 K 线
func (api *FutuAPI) GetCurKL(ctx context.Context, sec *Security, num int32, rehabType qotcommon.RehabType, klType qotcommon.KLType) (*CurKLine, error) {
	ch := make(qotgetkl.ResponseChan)
	if err := api.get(ProtoIDQotGetKL, &qotgetkl.Request{
		C2S: &qotgetkl.C2S{
			Security:  sec.pb(),
			ReqNum:    &num,
			RehabType: (*int32)(&rehabType),
			KlType:    (*int32)(&klType),
		},
	}, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return curKLineFromPB(resp.GetS2C()), result(resp)
	}
}

type CurKLine struct {
	Security *Security
	KLines   []*KLine
}

func curKLineFromPB(pb *qotgetkl.S2C) *CurKLine {
	if pb == nil {
		return nil
	}
	return &CurKLine{
		Security: securityFromPB(pb.GetSecurity()),
		KLines:   kLineListFromPB(pb.GetKlList()),
	}
}
