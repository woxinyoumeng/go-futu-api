package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetticker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetTicker = 3010 //Qot_GetTicker	获取逐笔
)

func (api *FutuAPI) GetRTTicker(ctx context.Context, sec *Security, num int32) (*RTTicker, error) {
	ch := make(getTickerChan)
	if err := api.get(ProtoIDQotGetTicker, &qotgetticker.Request{C2S: &qotgetticker.C2S{
		Security:  sec.pb(),
		MaxRetNum: &num,
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
		return rtTickerFromGetPB(resp.GetS2C()), result(resp)
	}
}

func rtTickerFromGetPB(pb *qotgetticker.S2C) *RTTicker {
	if pb == nil {
		return nil
	}
	return &RTTicker{
		Security: securityFromPB(pb.GetSecurity()),
		Tickers:  tickerListFromPB(pb.GetTickerList()),
	}
}

type getTickerChan chan *qotgetticker.Response

var _ protocol.RespChan = make(getTickerChan)

func (ch getTickerChan) Send(b []byte) error {
	var resp qotgetticker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch getTickerChan) Close() {
	close(ch)
}
