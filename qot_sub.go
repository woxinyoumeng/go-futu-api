package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotsub"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotSub = 3001 //Qot_Sub	订阅或者反订阅
)

// 订阅注册需要的实时信息，指定股票和订阅的数据类型即可。
// 香港市场（含正股、窝轮、牛熊、期权、期货）订阅，需要 LV1 及以上的权限，BMP 权限下不支持订阅。
func (api *FutuAPI) Subscribe(ctx context.Context, securities []*Security, subTypes []qotcommon.SubType,
	isRegPush bool, isFirstPush bool, isSubOrderBookDetail bool, isExtendedTime bool) error {
	return api.qotSub(ctx, &qotSubReq{
		IsSub:                true,
		Securities:           securities,
		SubTypes:             subTypes,
		IsRegPush:            isRegPush,
		IsFirstPush:          isFirstPush,
		IsSubOrderBookDetail: isSubOrderBookDetail,
		IsExtendedTime:       isExtendedTime,
	})
}

// 取消订阅
func (api *FutuAPI) Unsubscribe(ctx context.Context, securities []*Security, subTypes []qotcommon.SubType) error {
	return api.qotSub(ctx, &qotSubReq{
		Securities: securities,
		SubTypes:   subTypes,
	})
}

// 取消所有订阅
func (api *FutuAPI) UnsubscribeAll(ctx context.Context) error {
	return api.qotSub(ctx, &qotSubReq{
		IsUnsubAll: true,
	})
}

type qotSubReq struct {
	IsSub                bool
	Securities           []*Security
	SubTypes             []qotcommon.SubType
	RehabTypes           []qotcommon.RehabType
	IsRegPush            bool
	IsFirstPush          bool
	IsSubOrderBookDetail bool
	IsExtendedTime       bool
	IsUnsubAll           bool
}

func (req qotSubReq) pb() *qotsub.Request {
	r := qotsub.Request{
		C2S: &qotsub.C2S{
			IsSubOrUnSub:         &req.IsSub,
			IsRegOrUnRegPush:     &req.IsRegPush,
			IsFirstPush:          &req.IsFirstPush,
			IsUnsubAll:           &req.IsUnsubAll,
			IsSubOrderBookDetail: &req.IsSubOrderBookDetail,
			ExtendedTime:         &req.IsExtendedTime,
		},
	}
	if req.Securities != nil {
		r.C2S.SecurityList = make([]*qotcommon.Security, len(req.Securities))
		for i, v := range req.Securities {
			r.C2S.SecurityList[i] = v.pb()
		}
	}
	if req.SubTypes != nil {
		r.C2S.SubTypeList = make([]int32, len(req.SubTypes))
		for i, v := range req.SubTypes {
			r.C2S.SubTypeList[i] = int32(v)
		}
	}
	if req.RehabTypes != nil {
		r.C2S.RegPushRehabTypeList = make([]int32, len(req.RehabTypes))
		for i, v := range req.RehabTypes {
			r.C2S.RegPushRehabTypeList[i] = int32(v)
		}
	}
	return &r
}

func (api *FutuAPI) qotSub(ctx context.Context, req *qotSubReq) error {
	ch := make(qotSubChan)
	if err := api.get(ProtoIDQotSub, req.pb(), ch); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		return result(resp)
	}
}

type qotSubChan chan *qotsub.Response

var _ protocol.RespChan = make(qotSubChan)

func (ch qotSubChan) Send(b []byte) error {
	var resp qotsub.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &resp
	return nil
}

func (ch qotSubChan) Close() {
	close(ch)
}
