package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatert"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateRT = 3009 //Qot_UpdateRT	推送分时
)

// 实时分时回调
func (api *FutuAPI) UpdateRT(ctx context.Context) (*UpdateRTChan, error) {
	ch := UpdateRTChan{
		RT:  make(chan *RTData),
		Err: make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateRT, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

type UpdateRTChan struct {
	RT  chan *RTData
	Err chan error
}

var _ protocol.RespChan = (*UpdateRTChan)(nil)

func (ch *UpdateRTChan) Send(b []byte) error {
	var resp qotupdatert.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.RT <- rtDataFromUpdatePB(resp.GetS2C())
	}
	return nil
}

func (ch *UpdateRTChan) Close() {
	close(ch.RT)
	close(ch.Err)
}

func rtDataFromUpdatePB(pb *qotupdatert.S2C) *RTData {
	if pb == nil {
		return nil
	}
	return &RTData{
		Security:   securityFromPB(pb.GetSecurity()),
		TimeShares: timeShareListFromPB(pb.GetRtList()),
	}
}
