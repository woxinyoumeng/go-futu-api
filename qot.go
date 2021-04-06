package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
	"github.com/hurisheng/go-futu-api/pb/qotgetsubinfo"
	"github.com/hurisheng/go-futu-api/pb/qotsub"
	"github.com/hurisheng/go-futu-api/pb/qotupdateticker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

// 订阅注册参数
type SubscribeReq struct {
	Securities           []*Security
	SubTypes             []qotcommon.SubType
	IsRegPush            bool
	IsFirstPush          bool
	IsSubOrderBookDetail bool
	IsExtendedTime       bool
}

// 订阅注册需要的实时信息，指定股票和订阅的数据类型即可。
// 香港市场（含正股、窝轮、牛熊、期权、期货）订阅，需要 LV1 及以上的权限，BMP 权限下不支持订阅。
func (api *FutuAPI) Subscribe(ctx context.Context, req *SubscribeReq) error {
	return api.qotSub(ctx, &qotSubReq{
		IsSub:                true,
		Securities:           req.Securities,
		SubTypes:             req.SubTypes,
		IsRegPush:            req.IsRegPush,
		IsFirstPush:          req.IsFirstPush,
		IsSubOrderBookDetail: req.IsSubOrderBookDetail,
		IsExtendedTime:       req.IsExtendedTime,
	})
}

// 取消订阅参数
type UnsubscribeReq struct {
	Securities []*Security
	SubTypes   []qotcommon.SubType
}

// 取消订阅
func (api *FutuAPI) Unsubscribe(ctx context.Context, req *UnsubscribeReq) error {
	return api.qotSub(ctx, &qotSubReq{
		Securities: req.Securities,
		SubTypes:   req.SubTypes,
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
	out := make(chan *qotsub.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return err
	}
	if err := api.send(ProtoIDQotSub, req.pb(), ch); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ErrInterrupted
	case resp, ok := <-out:
		if !ok {
			return ErrChannelClosed
		}
		return result(resp)
	}
}

type Subscription struct {
	ConnSubInfos   []*ConnSubInfo
	TotalUsedQuota int32
	RemainQuota    int32
}

func subscriptionFromPB(pb *qotgetsubinfo.S2C) *Subscription {
	if pb == nil {
		return nil
	}
	sub := Subscription{
		TotalUsedQuota: pb.GetTotalUsedQuota(),
		RemainQuota:    pb.GetRemainQuota(),
	}
	if list := pb.GetConnSubInfoList(); list != nil {
		sub.ConnSubInfos = make([]*ConnSubInfo, len(list))
		for i, v := range list {
			sub.ConnSubInfos[i] = connSubInfoFromPB(v)
		}
	}
	return &sub
}

// 获取订阅信息
func (api *FutuAPI) QuerySubscription(ctx context.Context, isAll bool) (*Subscription, error) {
	out := make(chan *qotgetsubinfo.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return nil, err
	}
	if err := api.send(ProtoIDQotGetSubInfo, &qotgetsubinfo.Request{C2S: &qotgetsubinfo.C2S{
		IsReqAllConn: &isAll,
	}}, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-out:
		if !ok {
			return nil, ErrChannelClosed
		}
		return subscriptionFromPB(resp.GetS2C()), result(resp)
	}
}

type BasicQot struct {
	Security        *Security                //股票
	IsSuspended     bool                     //是否停牌
	ListTime        string                   //上市日期字符串
	PriceSpread     float64                  //价差
	UpdateTime      string                   //最新价的更新时间字符串，对其他字段不适用
	HighPrice       float64                  //最高价
	OpenPrice       float64                  //开盘价
	LowPrice        float64                  //最低价
	CurPrice        float64                  //最新价
	LastClosePrice  float64                  //昨收价
	Volume          int64                    //成交量
	Turnover        float64                  //成交额
	TurnoverRate    float64                  //换手率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Amplitude       float64                  //振幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	DarkStatus      qotcommon.DarkStatus     //DarkStatus, 暗盘交易状态
	OptionExData    *OptionBasicQotExData    //期权特有字段
	ListTimestamp   float64                  //上市日期时间戳
	UpdateTimestamp float64                  //最新价的更新时间戳，对其他字段不适用
	PreMarket       *PreAfterMarketData      //盘前数据
	AfterMarket     *PreAfterMarketData      //盘后数据
	SecStatus       qotcommon.SecurityStatus //SecurityStatus, 股票状态
	FutureExData    *FutureBasicQotExData    //期货特有字段
}

func basicQotFromPB(pb *qotcommon.BasicQot) *BasicQot {
	if pb == nil {
		return nil
	}
	return &BasicQot{
		Security:        securityFromPB(pb.GetSecurity()),
		IsSuspended:     pb.GetIsSuspended(),
		ListTime:        pb.GetListTime(),
		PriceSpread:     pb.GetPriceSpread(),
		UpdateTime:      pb.GetUpdateTime(),
		HighPrice:       pb.GetHighPrice(),
		OpenPrice:       pb.GetOpenPrice(),
		LowPrice:        pb.GetLowPrice(),
		CurPrice:        pb.GetCurPrice(),
		LastClosePrice:  pb.GetLastClosePrice(),
		Volume:          pb.GetVolume(),
		Turnover:        pb.GetTurnover(),
		TurnoverRate:    pb.GetTurnoverRate(),
		Amplitude:       pb.GetAmplitude(),
		DarkStatus:      qotcommon.DarkStatus(pb.GetDarkStatus()),
		OptionExData:    optionBasicQotExDataFromPB(pb.GetOptionExData()),
		ListTimestamp:   pb.GetListTimestamp(),
		UpdateTimestamp: pb.GetUpdateTimestamp(),
		PreMarket:       preAfterMarketDataFromPB(pb.GetPreMarket()),
		AfterMarket:     preAfterMarketDataFromPB(pb.GetAfterMarket()),
		SecStatus:       qotcommon.SecurityStatus(pb.GetSecStatus()),
		FutureExData:    futureBasicQotExDataFromPB(pb.GetFutureExData()),
	}
}

// 获取股票基本行情
func (api *FutuAPI) QotGetBasicQot(ctx context.Context, securities []*Security) ([]*BasicQot, error) {
	req := qotgetbasicqot.Request{
		C2S: &qotgetbasicqot.C2S{},
	}
	if securities != nil {
		req.C2S.SecurityList = make([]*qotcommon.Security, len(securities))
		for i, v := range securities {
			req.C2S.SecurityList[i] = v.pb()
		}
	}
	out := make(chan *qotgetbasicqot.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return nil, err
	}
	if err := api.send(ProtoIDQotGetBasicQot, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-out:
		if !ok {
			return nil, ErrChannelClosed
		}
		var basic []*BasicQot
		if list := resp.GetS2C().GetBasicQotList(); list != nil {
			basic = make([]*BasicQot, len(list))
			for i, v := range list {
				basic[i] = basicQotFromPB(v)
			}
		}
		return basic, nil
	}
}

// // QotUpdateBasicQot 推送股票基本报价
// func (api *FutuAPI) QotUpdateBasicQot() (<-chan *Qot_UpdateBasicQot.Response, error) {
// 	out := make(chan *Qot_UpdateBasicQot.Response)
// 	if err := api.subscribe(ProtoIDQotUpdateBasicQot, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdateBasicQot error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetKL 推送股票基本报价
// func (api *FutuAPI) QotGetKL(req *Qot_GetKL.Request) (<-chan *Qot_GetKL.Response, error) {
// 	out := make(chan *Qot_GetKL.Response)
// 	if err := api.send(ProtoIDQotGetKL, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetKL error: %w", err)
// 	}
// 	return out, nil
// }

// // QotUpdateKL 推送股票基本报价
// func (api *FutuAPI) QotUpdateKL() (<-chan *Qot_UpdateKL.Response, error) {
// 	out := make(chan *Qot_UpdateKL.Response)
// 	if err := api.subscribe(ProtoIDQotUpdateKL, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdateKL error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetRT 获取分时
// func (api *FutuAPI) QotGetRT(req *Qot_GetRT.Request) (<-chan *Qot_GetRT.Response, error) {
// 	out := make(chan *Qot_GetRT.Response)
// 	if err := api.send(ProtoIDQotGetRT, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetRT error: %w", err)
// 	}
// 	return out, nil
// }

// // QotUpdateRT 推送分时
// func (api *FutuAPI) QotUpdateRT() (<-chan *Qot_UpdateRT.Response, error) {
// 	out := make(chan *Qot_UpdateRT.Response)
// 	if err := api.subscribe(ProtoIDQotUpdateRT, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdateRT error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetTicker 获取逐笔
// func (api *FutuAPI) QotGetTicker(req *Qot_GetTicker.Request) (<-chan *Qot_GetTicker.Response, error) {
// 	out := make(chan *Qot_GetTicker.Response)
// 	if err := api.send(ProtoIDQotGetTicker, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetTicker error: %w", err)
// 	}
// 	return out, nil
// }

type TickerItem struct {
	Time         string                    //时间字符串
	Sequence     int64                     // 唯一标识
	Dir          qotcommon.TickerDirection //TickerDirection, 买卖方向
	Price        float64                   //价格
	Volume       int64                     //成交量
	Turnover     float64                   //成交额
	RecvTime     float64                   //收到推送数据的本地时间戳，用于定位延迟
	Type         qotcommon.TickerType      //TickerType, 逐笔类型
	TypeSign     int32                     //逐笔类型符号
	PushDataType qotcommon.PushDataType    //用于区分推送情况，仅推送时有该字段
	Timestamp    float64                   //时间戳
}

func tickerItemFromPB(pb *qotcommon.Ticker) *TickerItem {
	if pb == nil {
		return nil
	}
	return &TickerItem{
		Time:         pb.GetTime(),
		Sequence:     pb.GetSequence(),
		Dir:          qotcommon.TickerDirection(pb.GetDir()),
		Price:        pb.GetPrice(),
		Volume:       pb.GetVolume(),
		Turnover:     pb.GetTurnover(),
		RecvTime:     pb.GetRecvTime(),
		Type:         qotcommon.TickerType(pb.GetType()),
		TypeSign:     pb.GetTypeSign(),
		PushDataType: qotcommon.PushDataType(pb.GetPushDataType()),
		Timestamp:    pb.GetTimestamp(),
	}
}

type Ticker struct {
	Security *Security
	Items    []*TickerItem
}

func tickerFromPB(pb *qotupdateticker.S2C) *Ticker {
	if pb == nil {
		return nil
	}
	t := Ticker{
		Security: securityFromPB(pb.GetSecurity()),
	}
	if list := pb.GetTickerList(); list != nil {
		t.Items = make([]*TickerItem, len(list))
		for i, v := range list {
			t.Items[i] = tickerItemFromPB(v)
		}
	}
	return &t
}

type updateTickerChan struct {
	ticker chan *Ticker
	err    chan error
}

var _ protocol.RespChan = (*updateTickerChan)(nil)

func newUpdateTickerChan(tCh chan *Ticker, eCh chan error) *updateTickerChan {
	return &updateTickerChan{
		ticker: tCh,
		err:    eCh,
	}
}

func (ch *updateTickerChan) Send(b []byte) error {
	var resp qotupdateticker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	if err := result(&resp); err != nil {
		ch.err <- err
	} else {
		ch.ticker <- tickerFromPB(resp.GetS2C())
	}
	return nil
}
func (ch *updateTickerChan) Close() {
	close(ch.ticker)
	close(ch.err)
}

// 实时逐笔回调，异步处理已订阅股票的实时逐笔推送
func (api *FutuAPI) QotUpdateTicker(ctx context.Context) (<-chan *Ticker, <-chan error, error) {
	tCh := make(chan *Ticker)
	eCh := make(chan error)
	if err := api.subscribe(ProtoIDQotUpdateTicker, newUpdateTickerChan(tCh, eCh)); err != nil {
		return nil, nil, err
	}
	return tCh, eCh, nil
}

// // QotGetOrderBook 获取买卖盘
// func (api *FutuAPI) QotGetOrderBook(req *Qot_GetOrderBook.Request) (<-chan *Qot_GetOrderBook.Response, error) {
// 	out := make(chan *Qot_GetOrderBook.Response)
// 	if err := api.send(ProtoIDQotGetOrderBook, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetOrderBook error: %w", err)
// 	}
// 	return out, nil
// }

// // QotUpdateOrderBook 推送买卖盘
// func (api *FutuAPI) QotUpdateOrderBook() (<-chan *Qot_UpdateOrderBook.Response, error) {
// 	out := make(chan *Qot_UpdateOrderBook.Response)
// 	if err := api.subscribe(ProtoIDQotUpdateOrderBook, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdateOrderBook error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetBroker 获取经纪队列
// func (api *FutuAPI) QotGetBroker(req *Qot_GetBroker.Request) (<-chan *Qot_GetBroker.Response, error) {
// 	out := make(chan *Qot_GetBroker.Response)
// 	if err := api.send(ProtoIDQotGetBroker, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetBroker error: %w", err)
// 	}
// 	return out, nil
// }

// // QotUpdateBroker 推送经纪队列
// func (api *FutuAPI) QotUpdateBroker() (<-chan *Qot_UpdateBroker.Response, error) {
// 	out := make(chan *Qot_UpdateBroker.Response)
// 	if err := api.subscribe(ProtoIDQotUpdateBroker, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdateBroker error: %w", err)
// 	}
// 	return out, nil
// }

// // QotRequestRehab 在线获取单只股票复权信息
// func (api *FutuAPI) QotRequestRehab(req *Qot_RequestRehab.Request) (<-chan *Qot_RequestRehab.Response, error) {
// 	out := make(chan *Qot_RequestRehab.Response)
// 	if err := api.send(ProtoIDQotRequestRehab, req, out); err != nil {
// 		return nil, fmt.Errorf("QotRequestRehab error: %w", err)
// 	}
// 	return out, nil
// }

// // QotRequestHistoryKL 在线获取单只股票一段历史K线
// func (api *FutuAPI) QotRequestHistoryKL(req *Qot_RequestHistoryKL.Request) (<-chan *Qot_RequestHistoryKL.Response, error) {
// 	out := make(chan *Qot_RequestHistoryKL.Response)
// 	if err := api.send(ProtoIDQotRequestHistoryKL, req, out); err != nil {
// 		return nil, fmt.Errorf("QotRequestHistoryKL error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetStaticInfo 获取股票静态信息
// func (api *FutuAPI) QotGetStaticInfo(req *Qot_GetStaticInfo.Request) (<-chan *Qot_GetStaticInfo.Response, error) {
// 	out := make(chan *Qot_GetStaticInfo.Response)
// 	if err := api.send(ProtoIDQotGetStaticInfo, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetStaticInfo error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetSecuritySnapshot 获取股票快照
// func (api *FutuAPI) QotGetSecuritySnapshot(req *Qot_GetSecuritySnapshot.Request) (<-chan *Qot_GetSecuritySnapshot.Response, error) {
// 	out := make(chan *Qot_GetSecuritySnapshot.Response)
// 	if err := api.subscribe(ProtoIDQotGetSecuritySnapshot, out); err != nil {
// 		return nil, fmt.Errorf("QotGetSecuritySnapshot error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetPlateSet 获取板块集合下的板块
// func (api *FutuAPI) QotGetPlateSet(req *Qot_GetPlateSet.Request) (<-chan *Qot_GetPlateSet.Response, error) {
// 	out := make(chan *Qot_GetPlateSet.Response)
// 	if err := api.send(ProtoIDQotGetPlateSet, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetPlateSet error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetPlateSecurity 获取板块下的股票
// func (api *FutuAPI) QotGetPlateSecurity(req *Qot_GetPlateSecurity.Request) (<-chan *Qot_GetPlateSecurity.Response, error) {
// 	out := make(chan *Qot_GetPlateSecurity.Response)
// 	if err := api.send(ProtoIDQotGetPlateSecurity, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetPlateSecurity error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetReference 获取正股相关股票
// func (api *FutuAPI) QotGetReference(req *Qot_GetReference.Request) (<-chan *Qot_GetReference.Response, error) {
// 	out := make(chan *Qot_GetReference.Response)
// 	if err := api.send(ProtoIDQotGetReference, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetReference error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetOwnerPlate 获取股票所属板块
// func (api *FutuAPI) QotGetOwnerPlate(req *Qot_GetOwnerPlate.Request) (<-chan *Qot_GetOwnerPlate.Response, error) {
// 	out := make(chan *Qot_GetOwnerPlate.Response)
// 	if err := api.send(ProtoIDQotGetOwnerPlate, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetOwnerPlate error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetHoldingChangeList 获取持股变化列表
// func (api *FutuAPI) QotGetHoldingChangeList(req *Qot_GetHoldingChangeList.Request) (<-chan *Qot_GetHoldingChangeList.Response, error) {
// 	out := make(chan *Qot_GetHoldingChangeList.Response)
// 	if err := api.send(ProtoIDQotGetHoldingChangeList, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetHoldingChangeList error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetOptionChain 获取期权链
// func (api *FutuAPI) QotGetOptionChain(req *Qot_GetOptionChain.Request) (<-chan *Qot_GetOptionChain.Response, error) {
// 	out := make(chan *Qot_GetOptionChain.Response)
// 	if err := api.send(ProtoIDQotGetOptionChain, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetOptionChain error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetWarrant 获取窝轮
// func (api *FutuAPI) QotGetWarrant(req *Qot_GetWarrant.Request) (<-chan *Qot_GetWarrant.Response, error) {
// 	out := make(chan *Qot_GetWarrant.Response)
// 	if err := api.send(ProtoIDQotGetWarrant, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetWarrant error: %w", err)
// 	}
// 	return out, nil
// }

// // QotRequestHistoryKLQuota 拉取历史K线已经用掉的额度
// func (api *FutuAPI) QotRequestHistoryKLQuota(req *Qot_RequestHistoryKLQuota.Request) (<-chan *Qot_RequestHistoryKLQuota.Response, error) {
// 	out := make(chan *Qot_RequestHistoryKLQuota.Response)
// 	if err := api.send(ProtoIDQotRequestHistoryKLQuota, req, out); err != nil {
// 		return nil, fmt.Errorf("QotRequestHistoryKLQuota error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetCapitalFlow 获取资金流向
// func (api *FutuAPI) QotGetCapitalFlow(req *Qot_GetCapitalFlow.Request) (<-chan *Qot_GetCapitalFlow.Response, error) {
// 	out := make(chan *Qot_GetCapitalFlow.Response)
// 	if err := api.send(ProtoIDQotGetCapitalFlow, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetCapitalFlow error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetCapitalDistribution 获取资金分布
// func (api *FutuAPI) QotGetCapitalDistribution(req *Qot_GetCapitalDistribution.Request) (<-chan *Qot_GetCapitalDistribution.Response, error) {
// 	out := make(chan *Qot_GetCapitalDistribution.Response)
// 	if err := api.send(ProtoIDQotGetCapitalDistribution, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetCapitalDistribution error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetUserSecurity 获取自选股分组下的股票
// func (api *FutuAPI) QotGetUserSecurity(req *Qot_GetUserSecurity.Request) (<-chan *Qot_GetUserSecurity.Response, error) {
// 	out := make(chan *Qot_GetUserSecurity.Response)
// 	if err := api.send(ProtoIDQotGetUserSecurity, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetUserSecurity error: %w", err)
// 	}
// 	return out, nil
// }

// // QotModifyUserSecurity 修改自选股分组下的股票
// func (api *FutuAPI) QotModifyUserSecurity(req *Qot_ModifyUserSecurity.Request) (<-chan *Qot_ModifyUserSecurity.Response, error) {
// 	out := make(chan *Qot_ModifyUserSecurity.Response)
// 	if err := api.send(ProtoIDQotModifyUserSecurity, req, out); err != nil {
// 		return nil, fmt.Errorf("QotModifyUserSecurity error: %w", err)
// 	}
// 	return out, nil
// }

// // QotStockFilter 获取条件选股
// func (api *FutuAPI) QotStockFilter(req *Qot_StockFilter.Request) (<-chan *Qot_StockFilter.Response, error) {
// 	out := make(chan *Qot_StockFilter.Response)
// 	if err := api.send(ProtoIDQotStockFilter, req, out); err != nil {
// 		return nil, fmt.Errorf("QotStockFilter error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetCodeChange 获取股票代码变更信息
// func (api *FutuAPI) QotGetCodeChange(req *Qot_GetCodeChange.Request) (<-chan *Qot_GetCodeChange.Response, error) {
// 	out := make(chan *Qot_GetCodeChange.Response)
// 	if err := api.send(ProtoIDQotGetCodeChange, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetCodeChange error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetIpoList 获取IPO信息
// func (api *FutuAPI) QotGetIpoList(req *Qot_GetIpoList.Request) (<-chan *Qot_GetIpoList.Response, error) {
// 	out := make(chan *Qot_GetIpoList.Response)
// 	if err := api.send(ProtoIDQotGetIpoList, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetIpoList error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetFutureInfo 获取期货合约资料
// func (api *FutuAPI) QotGetFutureInfo(req *Qot_GetFutureInfo.Request) (<-chan *Qot_GetFutureInfo.Response, error) {
// 	out := make(chan *Qot_GetFutureInfo.Response)
// 	if err := api.send(ProtoIDQotGetFutureInfo, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetFutureInfo error: %w", err)
// 	}
// 	return out, nil
// }

// // QotRequestTradeDate 在线请求交易日
// func (api *FutuAPI) QotRequestTradeDate(req *Qot_RequestTradeDate.Request) (<-chan *Qot_RequestTradeDate.Response, error) {
// 	out := make(chan *Qot_RequestTradeDate.Response)
// 	if err := api.send(ProtoIDQotRequestTradeDate, req, out); err != nil {
// 		return nil, fmt.Errorf("QotRequestTradeDate error: %w", err)
// 	}
// 	return out, nil
// }

// // QotSetPriceReminder 设置到价提醒
// func (api *FutuAPI) QotSetPriceReminder(req *Qot_SetPriceReminder.Request) (<-chan *Qot_SetPriceReminder.Response, error) {
// 	out := make(chan *Qot_SetPriceReminder.Response)
// 	if err := api.send(ProtoIDQotSetPriceReminder, req, out); err != nil {
// 		return nil, fmt.Errorf("QotSetPriceReminder error: %w", err)
// 	}
// 	return out, nil
// }

// // QotGetPriceReminder 获取到价提醒
// func (api *FutuAPI) QotGetPriceReminder(req *Qot_GetPriceReminder.Request) (<-chan *Qot_GetPriceReminder.Response, error) {
// 	out := make(chan *Qot_GetPriceReminder.Response)
// 	if err := api.send(ProtoIDQotGetPriceReminder, req, out); err != nil {
// 		return nil, fmt.Errorf("QotGetPriceReminder error: %w", err)
// 	}
// 	return out, nil
// }

// // QotUpdatePriceReminder 到价提醒通知
// func (api *FutuAPI) QotUpdatePriceReminder() (<-chan *Qot_UpdatePriceReminder.Response, error) {
// 	out := make(chan *Qot_UpdatePriceReminder.Response)
// 	if err := api.subscribe(ProtoIDQotUpdatePriceReminder, out); err != nil {
// 		return nil, fmt.Errorf("QotUpdatePriceReminder error: %w", err)
// 	}
// 	return out, nil
// }
