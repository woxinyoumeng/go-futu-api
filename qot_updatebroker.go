package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatebroker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateBroker = 3015 //Qot_UpdateBroker	推送经纪队列
)

// 实时经纪队列回调
func (api *FutuAPI) UpdateBroker(ctx context.Context) (*UpdateBrokerChan, error) {
	ch := UpdateBrokerChan{
		BrokerQueue: make(chan *BrokerQueue),
		Err:         make(chan error),
	}
	if err := api.update(ProtoIDQotUpdateBroker, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

type UpdateBrokerChan struct {
	BrokerQueue chan *BrokerQueue
	Err         chan error
}

var _ protocol.RespChan = (*UpdateBrokerChan)(nil)

func (ch *UpdateBrokerChan) Send(b []byte) error {
	var resp qotupdatebroker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.Err <- err
	} else {
		ch.BrokerQueue <- brokerQueueFromUpdatePB(resp.GetS2C())
	}
	return nil
}

func (ch *UpdateBrokerChan) Close() {
	close(ch.BrokerQueue)
	close(ch.Err)
}

func brokerQueueFromUpdatePB(pb *qotupdatebroker.S2C) *BrokerQueue {
	if pb == nil {
		return nil
	}
	return &BrokerQueue{
		Security: securityFromPB(pb.GetSecurity()),
		Asks:     brokerListFromPB(pb.GetBrokerAskList()),
		Bids:     brokerListFromPB(pb.GetBrokerBidList()),
	}
}
