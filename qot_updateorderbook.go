package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdateorderbook"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateOrderBook = 3013 //Qot_UpdateOrderBook	推送买卖盘
)

// 实时摆盘回调
func (api *FutuAPI) UpdateOrderBook(ctx context.Context) (*UpdateOrderBookChan, error) {
	ch := UpdateOrderBookChan{
		OrderBook: make(chan *RTOrderBook),
		Err:       make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateOrderBook, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

type UpdateOrderBookChan struct {
	OrderBook chan *RTOrderBook
	Err       chan error
}

var _ protocol.RespChan = (*UpdateOrderBookChan)(nil)

func (ch *UpdateOrderBookChan) Send(b []byte) error {
	var resp qotupdateorderbook.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.OrderBook <- rtOrderBookFromUpdatePB(resp.GetS2C())
	}
	return nil
}

func (ch *UpdateOrderBookChan) Close() {
	close(ch.OrderBook)
	close(ch.Err)
}

func rtOrderBookFromUpdatePB(pb *qotupdateorderbook.S2C) *RTOrderBook {
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
