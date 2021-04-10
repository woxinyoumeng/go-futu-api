package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotupdatekl"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateKL = 3007 //Qot_UpdateKL	推送 K 线
)

// 实时 K 线回调
func (api *FutuAPI) UpdateKL(ctx context.Context) (*UpdateKLChan, error) {
	ch := UpdateKLChan{
		KLine: make(chan *UpdateKL),
		Err:   make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateKL, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// 实时 K 线推送
type UpdateKL struct {
	RehabType qotcommon.RehabType //*Qot_Common.RehabType,复权类型
	KLType    qotcommon.KLType    //*Qot_Common.KLType,K 线类型
	Security  *Security           //*股票
	KLines    []*KLine
}

func updateKLFromPB(pb *qotupdatekl.S2C) *UpdateKL {
	if pb == nil {
		return nil
	}
	return &UpdateKL{
		RehabType: qotcommon.RehabType(pb.GetRehabType()),
		KLType:    qotcommon.KLType(pb.GetKlType()),
		Security:  securityFromPB(pb.GetSecurity()),
		KLines:    kLineListFromPB(pb.GetKlList()),
	}
}

type UpdateKLChan struct {
	KLine chan *UpdateKL
	Err   chan error
}

var _ protocol.RespChan = (*UpdateKLChan)(nil)

func (ch *UpdateKLChan) Send(b []byte) error {
	var resp qotupdatekl.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.KLine <- updateKLFromPB(resp.GetS2C())
	}
	return nil
}

func (ch *UpdateKLChan) Close() {
	close(ch.KLine)
	close(ch.Err)
}
