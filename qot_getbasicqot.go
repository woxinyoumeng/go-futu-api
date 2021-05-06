package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
)

const (
	ProtoIDQotGetBasicQot = 3004 //Qot_GetBasicQot	获取股票基本报价
)

// 获取股票基本行情
func (api *FutuAPI) GetStockQuote(ctx context.Context, securities []*Security) ([]*BasicQot, error) {
	ch := make(qotgetbasicqot.ResponseChan)
	if err := api.get(ProtoIDQotGetBasicQot, &qotgetbasicqot.Request{
		C2S: &qotgetbasicqot.C2S{
			SecurityList: securityList(securities).pb(),
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
		return basicQotListFromPB(resp.GetS2C().GetBasicQotList()), nil
	}
}
