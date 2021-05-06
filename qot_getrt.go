package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetrt"
)

const (
	ProtoIDQotGetRT = 3008 //Qot_GetRT	获取分时
)

// 获取实时分时
func (api *FutuAPI) GetRTData(ctx context.Context, sec *Security) (*RTData, error) {
	ch := make(qotgetrt.ResponseChan)
	if err := api.get(ProtoIDQotGetRT, &qotgetrt.Request{C2S: &qotgetrt.C2S{
		Security: sec.pb(),
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
		return rtDataFromGetPB(resp.GetS2C()), result(resp)
	}
}

func rtDataFromGetPB(pb *qotgetrt.S2C) *RTData {
	if pb == nil {
		return nil
	}
	return &RTData{
		Security:   securityFromPB(pb.GetSecurity()),
		TimeShares: timeShareListFromPB(pb.GetRtList()),
	}
}
