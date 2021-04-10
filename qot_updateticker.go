package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdateticker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateTicker = 3011 //Qot_UpdateTicker	推送逐笔
)

// 实时逐笔回调，异步处理已订阅股票的实时逐笔推送
func (api *FutuAPI) UpdateTicker(ctx context.Context) (*UpdateTickerChan, error) {
	ch := UpdateTickerChan{
		Ticker: make(chan *RTTicker),
		Err:    make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateTicker, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

type UpdateTickerChan struct {
	Ticker chan *RTTicker
	Err    chan error
}

var _ protocol.RespChan = (*UpdateTickerChan)(nil)

func (ch *UpdateTickerChan) Send(b []byte) error {
	var resp qotupdateticker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.Ticker <- rtTickerFromUpdatePB(resp.GetS2C())
	}
	return nil
}
func (ch *UpdateTickerChan) Close() {
	close(ch.Ticker)
	close(ch.Err)
}

func rtTickerFromUpdatePB(pb *qotupdateticker.S2C) *RTTicker {
	if pb == nil {
		return nil
	}
	return &RTTicker{
		Security: securityFromPB(pb.GetSecurity()),
		Tickers:  tickerListFromPB(pb.GetTickerList()),
	}
}
