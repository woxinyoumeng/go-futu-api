package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetmarketstate"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetMarketState = 3223 //Qot_GetMarketState	获取指定品种的市场状态
)

//获取标的市场状态
func (api *FutuAPI) GetMarketState(ctx context.Context, securities []*Security) ([]*MarketInfo, error) {
	ch := make(marketStateChan)
	if err := api.get(ProtoIDQotGetMarketState, &qotgetmarketstate.Request{
		C2S: &qotgetmarketstate.C2S{
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
		return marketInfoListFromPB(resp.GetS2C().GetMarketInfoList()), result(resp)
	}
}

type marketStateChan chan *qotgetmarketstate.Response

var _ protocol.RespChan = make(marketStateChan)

func (ch marketStateChan) Send(b []byte) error {
	var resp qotgetmarketstate.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch marketStateChan) Close() {
	close(ch)
}

func marketInfoListFromPB(pb []*qotgetmarketstate.MarketInfo) []*MarketInfo {
	if pb == nil {
		return nil
	}
	m := make([]*MarketInfo, len(pb))
	for i, v := range pb {
		m[i] = marketInfoFromPB(v)
	}
	return m
}

type MarketInfo struct {
	Security    *Security                //股票代码
	Name        string                   // 股票名称
	MarketState qotcommon.QotMarketState //Qot_Common.QotMarketState，市场状态
}

func marketInfoFromPB(pb *qotgetmarketstate.MarketInfo) *MarketInfo {
	if pb == nil {
		return nil
	}
	return &MarketInfo{
		Security:    securityFromPB(pb.GetSecurity()),
		Name:        pb.GetName(),
		MarketState: qotcommon.QotMarketState(pb.GetMarketState()),
	}
}
