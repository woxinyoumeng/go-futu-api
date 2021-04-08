package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetBasicQot = 3004 //Qot_GetBasicQot	获取股票基本报价
)

// 获取股票基本行情
func (api *FutuAPI) GetStockQuote(ctx context.Context, securities []*Security) ([]*BasicQot, error) {
	req := qotgetbasicqot.Request{
		C2S: &qotgetbasicqot.C2S{},
	}
	if securities != nil {
		req.C2S.SecurityList = make([]*qotcommon.Security, len(securities))
		for i, v := range securities {
			req.C2S.SecurityList[i] = v.pb()
		}
	}
	ch := make(qotGetBasicQotChan)
	if err := api.get(ProtoIDQotGetBasicQot, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		var basic []*BasicQot
		if list := resp.GetS2C().GetBasicQotList(); list != nil {
			basic = make([]*BasicQot, len(list))
			for i, v := range list {
				basic[i] = basicQotFromPB(v)
			}
		}
		return basic, nil
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
