package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetorderbook"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetOrderBook = 3012 //Qot_GetOrderBook	获取买卖盘
)

// 获取实时摆盘
func (api *FutuAPI) GetOrderBook(ctx context.Context, sec *Security, num int32) (*RTOrderBook, error) {
	ch := make(getOrderBookChan)
	if err := api.get(ProtoIDQotGetOrderBook, &qotgetorderbook.Request{
		C2S: &qotgetorderbook.C2S{
			Security: sec.pb(),
			Num:      &num,
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
		return rtOrderBookFromGetPB(resp.GetS2C()), result(resp)
	}
}

func rtOrderBookFromGetPB(pb *qotgetorderbook.S2C) *RTOrderBook {
	if pb == nil {
		return nil
	}
	return &RTOrderBook{
		Security:                securityFromPB(pb.GetSecurity()),
		Asks:                    orderBookListFromPB(pb.GetOrderBookAskList()),
		Bids:                    orderBookListFromPB(pb.GetOrderBookBidList()),
		SvrRecvTimeBid:          pb.GetSvrRecvTimeBid(),
		SvrRecvTimeBidTimestamp: pb.GetSvrRecvTimeBidTimestamp(),
		SvrRecvTimeAsk:          pb.GetSvrRecvTimeAsk(),
		SvrRecvTimeAskTimestamp: pb.GetSvrRecvTimeAskTimestamp(),
	}
}

type getOrderBookChan chan *qotgetorderbook.Response

var _ protocol.RespChan = make(getOrderBookChan)

func (ch getOrderBookChan) Send(b []byte) error {
	var resp qotgetorderbook.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch getOrderBookChan) Close() {
	close(ch)
}
