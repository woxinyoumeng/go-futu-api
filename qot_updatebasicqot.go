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

func (api *FutuAPI) UpdateBasicQot(ctx context.Context) (<-chan []*BasicQot, <-chan error, error) {
	bCh := make(chan []*BasicQot)
	eCh := make(chan error)
	if err := api.update(ProtoIDQotUpdateBasicQot, &updateBasicQotChan{basic: bCh, err: eCh}); err != nil {
		return nil, nil, err
	}
	return bCh, eCh, nil
}

type updateBasicQotChan struct {
	basic chan []*BasicQot
	err   chan error
}

var _ protocol.RespChan = (*updateBasicQotChan)(nil)

func (ch *updateBasicQotChan) Send(b []byte) error {
	var resp qotupdatebasicqot.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.err <- err
	} else {
		var basic []*BasicQot
		if list := resp.GetS2C().GetBasicQotList(); list != nil {
			basic = make([]*BasicQot, len(list))
			for i, v := range list {
				basic[i] = basicQotFromPB(v)
			}
		}
		ch.basic <- basic
	}
	return nil
}

func (ch *updateBasicQotChan) Close() {
	close(ch.basic)
	close(ch.err)
}
