package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetrt"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetRT = 3008 //Qot_GetRT	获取分时
)

func (api *FutuAPI) GetRTData(ctx context.Context, sec *Security) (*RTData, error) {
	ch := make(rtDataChan)
	if err := api.get(ProtoIDQotGetRT, &qotgetrt.Request{C2S: &qotgetrt.C2S{
		Security: sec.pb(),
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
		return rtDataFromPB(resp.GetS2C()), result(resp)
	}
}

// 实时分时数据
type RTData struct {
	Security   *Security    //*股票
	TimeShares []*TimeShare //*分时数据结构体
}

func rtDataFromPB(pb *qotgetrt.S2C) *RTData {
	if pb == nil {
		return nil
	}
	rt := RTData{
		Security: securityFromPB(pb.GetSecurity()),
	}
	if list := pb.GetRtList(); list != nil {
		rt.TimeShares = make([]*TimeShare, len(list))
		for i, v := range list {
			rt.TimeShares[i] = timeShareFromPB(v)
		}
	}
	return &rt
}

type rtDataChan chan *qotgetrt.Response

var _ protocol.RespChan = make(rtDataChan)

func (ch rtDataChan) Send(b []byte) error {
	var resp qotgetrt.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch rtDataChan) Close() {
	close(ch)
}
