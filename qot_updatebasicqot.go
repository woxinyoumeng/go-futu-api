package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatebasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateBasicQot = 3005 //Qot_UpdateBasicQot	推送股票基本报价
)

// 实时报价回调
func (api *FutuAPI) UpdateBasicQot(ctx context.Context) (*UpdateBasicQotChan, error) {
	ch := UpdateBasicQotChan{
		BasicQot: make(chan []*BasicQot),
		Err:      make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateBasicQot, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

type UpdateBasicQotChan struct {
	BasicQot chan []*BasicQot
	Err      chan error
}

var _ protocol.RespChan = (*UpdateBasicQotChan)(nil)

func (ch *UpdateBasicQotChan) Send(b []byte) error {
	var resp qotupdatebasicqot.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.BasicQot <- basicQotListFromPB(resp.GetS2C().GetBasicQotList())
	}
	return nil
}

func (ch *UpdateBasicQotChan) Close() {
	close(ch.BasicQot)
	close(ch.Err)
}
