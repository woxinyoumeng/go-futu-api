package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetkl"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotGetKL = 3006 //Qot_GetKL	获取 K 线
)

func (api *FutuAPI) GetCurKL(ctx context.Context, sec *Security, num int32, rehabType qotcommon.RehabType, klType qotcommon.KLType) (*CurKLine, error) {
	ch := make(qotGetKLChan)
	if err := api.get(ProtoIDQotGetKL, &qotgetkl.Request{
		C2S: &qotgetkl.C2S{
			Security:  sec.pb(),
			ReqNum:    &num,
			RehabType: (*int32)(&rehabType),
			KlType:    (*int32)(&klType),
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
		return curKLineFromPB(resp.GetS2C()), result(resp)
	}
}

type CurKLine struct {
	Security *Security
	KLines   []*KLine
}

func curKLineFromPB(pb *qotgetkl.S2C) *CurKLine {
	if pb == nil {
		return nil
	}
	kl := CurKLine{
		Security: securityFromPB(pb.GetSecurity()),
	}
	if list := pb.GetKlList(); list != nil {
		kl.KLines = make([]*KLine, len(list))
		for i, v := range list {
			kl.KLines[i] = kLineFromPB(v)
		}
	}
	return &kl
}

type qotGetKLChan chan *qotgetkl.Response

var _ protocol.RespChan = make(qotGetKLChan)

func (ch qotGetKLChan) Send(b []byte) error {
	var resp qotgetkl.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch qotGetKLChan) Close() {
	close(ch)
}
