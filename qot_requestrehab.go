package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotrequestrehab"
)

const (
	ProtoIDQotRequestRehab = 3105 //Qot_RequestRehab	在线获取单只股票复权信息
)

// 获取复权因子
func (api *FutuAPI) GetRehab(ctx context.Context, security *Security) ([]*Rehab, error) {
	ch := make(qotrequestrehab.ResponseChan)
	if err := api.get(ProtoIDQotRequestRehab, &qotrequestrehab.Request{
		C2S: &qotrequestrehab.C2S{
			Security: security.pb(),
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
		return rehabListFromPB(resp.GetS2C().GetRehabList()), result(resp)
	}
}
