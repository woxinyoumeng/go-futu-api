package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetBasicQot = 3004 //Qot_GetBasicQot	获取股票基本报价
)

// 获取股票基本行情
func (api *FutuAPI) GetStockQuote(ctx context.Context, securities []*Security) ([]*BasicQot, error) {
	ch := make(qotGetBasicQotChan)
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

type qotGetBasicQotChan chan *qotgetbasicqot.Response

var _ protocol.RespChan = make(qotGetBasicQotChan)

func (ch qotGetBasicQotChan) Send(b []byte) error {
	var resp qotgetbasicqot.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	return nil
}

func (ch qotGetBasicQotChan) Close() {
	close(ch)
}
