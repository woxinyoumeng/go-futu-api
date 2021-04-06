package futuapi

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {
	api := NewFutuAPI()

	if err := api.Connect(context.Background(), ":11111"); err != nil {
		t.Error(err)
		return
	}
	resp, err := api.GetGlobalState(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Error(resp)
	// req := QotSubReq{
	// 	Securities: []*Security{{qotcommon.QotMarket_QotMarket_HK_Security, "00700"}},
	// 	IsSub:      true,
	// 	SubTypes:   []qotcommon.SubType{qotcommon.SubType_SubType_Ticker},
	// 	RehabTypes: []qotcommon.RehabType{qotcommon.RehabType_RehabType_Forward},
	// }
	// if err := api.QotSub(context.Background(), &req); err != nil {
	// 	t.Error(err)
	// }
	// ticker, err := api.QotUpdateTicker()
	// if err != nil {
	// 	t.Error(err)
	// }
	// for v := range ticker {
	// 	t.Error(v)
	// }
}
