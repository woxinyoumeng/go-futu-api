package futuapi

import (
	"context"
	"testing"

	"github.com/hurisheng/go-futu-api/protobuf/qotcommon"
	"github.com/hurisheng/go-futu-api/protobuf/qotsub"
)

func TestConnect(t *testing.T) {
	api, err := NewFutuAPI(&Config{
		Address:   ":11111",
		ClientVer: 100,
		ClientID:  "1",
	})
	if err != nil {
		t.Error(err)
	}
	var market int32 = int32(qotcommon.QotMarket_QotMarket_HK_Security)
	var code string = "0700"
	if re, err := api.QotSub(context.Background(), &qotsub.C2S{
		SecurityList: []*qotcommon.Security{
			{Market: &market, Code: &code},
		},
	}); err != nil {
		t.Error(err)
	} else {
		t.Error(re)
	}
	// ticker, err := api.QotUpdateTicker()
	// if err != nil {
	// 	t.Error(err)
	// }
	// for v := range ticker {
	// 	t.Error(v)
	// }
}
