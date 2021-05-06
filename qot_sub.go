package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotsub"
)

const (
	ProtoIDQotSub = 3001 //Qot_Sub	订阅或者反订阅
)

// 订阅注册需要的实时信息，指定股票和订阅的数据类型即可。
// 香港市场（含正股、窝轮、牛熊、期权、期货）订阅，需要 LV1 及以上的权限，BMP 权限下不支持订阅。
func (api *FutuAPI) Subscribe(ctx context.Context, securities []*Security, subTypes []qotcommon.SubType,
	isRegPush bool, isFirstPush bool, isSubOrderBookDetail bool, isExtendedTime bool) error {
	return api.qotSub(ctx, true, securities, subTypes, nil, isRegPush, isFirstPush, isSubOrderBookDetail, isExtendedTime, false)
}

// 取消订阅
func (api *FutuAPI) Unsubscribe(ctx context.Context, securities []*Security, subTypes []qotcommon.SubType) error {
	return api.qotSub(ctx, false, securities, subTypes, nil, false, false, false, false, false)
}

// 取消所有订阅
func (api *FutuAPI) UnsubscribeAll(ctx context.Context) error {
	return api.qotSub(ctx, false, nil, nil, nil, false, false, false, false, true)
}

func (api *FutuAPI) qotSub(ctx context.Context, isSub bool, securities []*Security, subTypes []qotcommon.SubType, rehabTypes []qotcommon.RehabType,
	isRegPush bool, isFirstPush bool, isSubOrderBookDetail bool, isExtendedTime bool, isUnsubAll bool) error {
	r := qotsub.Request{
		C2S: &qotsub.C2S{
			SecurityList:         securityList(securities).pb(),
			IsSubOrUnSub:         &isSub,
			IsRegOrUnRegPush:     &isRegPush,
			IsFirstPush:          &isFirstPush,
			IsUnsubAll:           &isUnsubAll,
			IsSubOrderBookDetail: &isSubOrderBookDetail,
			ExtendedTime:         &isExtendedTime,
		},
	}
	if subTypes != nil {
		r.C2S.SubTypeList = make([]int32, len(subTypes))
		for i, v := range subTypes {
			r.C2S.SubTypeList[i] = int32(v)
		}
	}
	if rehabTypes != nil {
		r.C2S.RegPushRehabTypeList = make([]int32, len(rehabTypes))
		for i, v := range rehabTypes {
			r.C2S.RegPushRehabTypeList[i] = int32(v)
		}
	}
	ch := make(qotsub.ResponseChan)
	if err := api.get(ProtoIDQotSub, &r, ch); err != nil {
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
