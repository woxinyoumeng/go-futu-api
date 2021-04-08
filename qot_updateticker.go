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
func (api *FutuAPI) UpdateTicker(ctx context.Context) (<-chan *Ticker, <-chan error, error) {
	tCh := make(chan *Ticker)
	eCh := make(chan error)
	if err := api.update(ProtoIDQotUpdateTicker, &updateTickerChan{ticker: tCh, err: eCh}); err != nil {
		return nil, nil, err
	}
	return tCh, eCh, nil
}

type updateTickerChan struct {
	ticker chan *Ticker
	err    chan error
}

var _ protocol.RespChan = (*updateTickerChan)(nil)

func (ch *updateTickerChan) Send(b []byte) error {
	var resp qotupdateticker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.err <- err
	} else {
		ch.ticker <- tickerFromPB(resp.GetS2C())
	}
	return nil
}
func (ch *updateTickerChan) Close() {
	close(ch.ticker)
	close(ch.err)
}
