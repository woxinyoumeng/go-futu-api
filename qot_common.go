package futuapi

import (
	"github.com/hurisheng/go-futu-api/pb/qotcommon"
)

// 证券标识
type Security struct {
	Market qotcommon.QotMarket //*QotMarket，股票市场
	Code   string              //*股票代码
}

func (s Security) pb() *qotcommon.Security {
	return &qotcommon.Security{
		Market: (*int32)(&s.Market),
		Code:   &s.Code,
	}
}

func securityFromPB(pb *qotcommon.Security) *Security {
	if pb == nil {
		return nil
	}
	return &Security{
		Market: qotcommon.QotMarket(pb.GetMarket()),
		Code:   pb.GetCode(),
	}
}

type securityList []*Security

func (s securityList) pb() []*qotcommon.Security {
	if s == nil {
		return nil
	}
	li := make([]*qotcommon.Security, len(s))
	for i, v := range s {
		li[i] = v.pb()
	}
	return li
}

// func securityListFromPB(pb []*qotcommon.Security) []*Security {
// 	if pb == nil {
// 		return nil
// 	}
// 	s := make(securityList, len(pb))
// 	for i, v := range pb {
// 		s[i] = securityFromPB(v)
// 	}
// 	return s
// }

// 基础报价的期权特有字段
type OptionBasicQotExData struct {
	StrikePrice          float64                   //*行权价
	ContractSize         int32                     //*每份合约数(整型数据)
	ContractSizeFloat    float64                   //每份合约数（浮点型数据）
	OpenInterest         int32                     //*未平仓合约数
	ImpliedVolatility    float64                   //*隐含波动率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Premium              float64                   //*溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Delta                float64                   //*希腊值 Delta
	Gamma                float64                   //*希腊值 Gamma
	Vega                 float64                   //*希腊值 Vega
	Theta                float64                   //*希腊值 Theta
	Rho                  float64                   //*希腊值 Rho
	NetOpenInterest      int32                     //净未平仓合约数，仅港股期权适用
	ExpiryDataDistance   int32                     //距离到期日天数，负数表示已过期
	ContractNominalValue float64                   //合约名义金额，仅港股期权适用
	OwnerLotMultiplier   float64                   //相等正股手数，指数期权无该字段，仅港股期权适用
	OptionAreaType       qotcommon.OptionAreaType  //OptionAreaType，期权类型（按行权时间）
	ContractMultiplier   float64                   //合约乘数
	IndexOptionType      qotcommon.IndexOptionType //IndexOptionType，指数期权类型
}

func optionBasicQotExDataFromPB(pb *qotcommon.OptionBasicQotExData) *OptionBasicQotExData {
	if pb == nil {
		return nil
	}
	return &OptionBasicQotExData{
		StrikePrice:          pb.GetStrikePrice(),
		ContractSize:         pb.GetContractSize(),
		ContractSizeFloat:    pb.GetContractSizeFloat(),
		OpenInterest:         pb.GetOpenInterest(),
		ImpliedVolatility:    pb.GetImpliedVolatility(),
		Premium:              pb.GetPremium(),
		Delta:                pb.GetDelta(),
		Gamma:                pb.GetGamma(),
		Vega:                 pb.GetVega(),
		Theta:                pb.GetTheta(),
		Rho:                  pb.GetRho(),
		NetOpenInterest:      pb.GetNetOpenInterest(),
		ExpiryDataDistance:   pb.GetExpiryDateDistance(),
		ContractNominalValue: pb.GetContractNominalValue(),
		OwnerLotMultiplier:   pb.GetOwnerLotMultiplier(),
		OptionAreaType:       qotcommon.OptionAreaType(pb.GetOptionAreaType()),
		ContractMultiplier:   pb.GetContractMultiplier(),
		IndexOptionType:      qotcommon.IndexOptionType(pb.GetIndexOptionType()),
	}
}

// 盘前盘后数据
type PreAfterMarketData struct {
	Price      float64 // 盘前或盘后## 价格
	HighPrice  float64 // 盘前或盘后## 最高价
	LowPrice   float64 // 盘前或盘后## 最低价
	Volume     int64   // 盘前或盘后## 成交量
	Turnover   float64 // 盘前或盘后## 成交额
	ChangeVal  float64 // 盘前或盘后## 涨跌额
	ChangeRate float64 // 盘前或盘后## 涨跌幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Amplitude  float64 // 盘前或盘后## 振幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
}

func preAfterMarketDataFromPB(pb *qotcommon.PreAfterMarketData) *PreAfterMarketData {
	if pb == nil {
		return nil
	}
	return &PreAfterMarketData{
		Price:      pb.GetPrice(),
		HighPrice:  pb.GetHighPrice(),
		LowPrice:   pb.GetLowPrice(),
		Volume:     pb.GetVolume(),
		Turnover:   pb.GetTurnover(),
		ChangeVal:  pb.GetChangeVal(),
		ChangeRate: pb.GetChangeRate(),
		Amplitude:  pb.GetAmplitude(),
	}
}

// 基础报价的期货特有字段
type FutureBasicQotExData struct {
	LastSettlePrice    float64 //*昨结
	Position           int32   //*持仓量
	PositionChange     int32   //*日增仓
	ExpiryDataDistance int32   //距离到期日天数
}

func futureBasicQotExDataFromPB(pb *qotcommon.FutureBasicQotExData) *FutureBasicQotExData {
	if pb == nil {
		return nil
	}
	return &FutureBasicQotExData{
		LastSettlePrice:    pb.GetLastSettlePrice(),
		Position:           pb.GetPosition(),
		PositionChange:     pb.GetPositionChange(),
		ExpiryDataDistance: pb.GetExpiryDateDistance(),
	}
}

// 基础报价
type BasicQot struct {
	Security        *Security                //*股票
	IsSuspended     bool                     //*是否停牌
	ListTime        string                   //*上市日期字符串
	PriceSpread     float64                  //*价差
	UpdateTime      string                   //*最新价的更新时间字符串，对其他字段不适用
	HighPrice       float64                  //*最高价
	OpenPrice       float64                  //*开盘价
	LowPrice        float64                  //*最低价
	CurPrice        float64                  //*最新价
	LastClosePrice  float64                  //*昨收价
	Volume          int64                    //*成交量
	Turnover        float64                  //*成交额
	TurnoverRate    float64                  //*换手率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Amplitude       float64                  //*振幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
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

func basicQotListFromPB(pb []*qotcommon.BasicQot) []*BasicQot {
	if pb == nil {
		return nil
	}
	bq := make([]*BasicQot, len(pb))
	for i, v := range pb {
		bq[i] = basicQotFromPB(v)
	}
	return bq
}

// 买卖档
type OrderBook struct {
	Price      float64            //*委托价格
	Volume     int64              //*委托数量
	OrderCount int32              //*委托订单个数
	Details    []*OrderBookDetail //订单信息，SF 行情特有
}

func orderBookFromPB(pb *qotcommon.OrderBook) *OrderBook {
	if pb == nil {
		return nil
	}
	item := OrderBook{
		Price:      pb.GetPrice(),
		Volume:     pb.GetVolume(),
		OrderCount: pb.GetOrederCount(),
	}
	if list := pb.GetDetailList(); list != nil {
		item.Details = make([]*OrderBookDetail, len(list))
		for i, v := range list {
			item.Details[i] = orderBookDetailFromPB(v)
		}
	}
	return &item
}

// 买卖档明细
type OrderBookDetail struct {
	OrderID int64 //交易所订单 ID，与交易接口返回的订单 ID 并不一样
	Volume  int64 //订单股数
}

func orderBookDetailFromPB(pb *qotcommon.OrderBookDetail) *OrderBookDetail {
	if pb == nil {
		return nil
	}
	return &OrderBookDetail{
		OrderID: pb.GetOrderID(),
		Volume:  pb.GetVolume(),
	}
}

func orderBookListFromPB(pb []*qotcommon.OrderBook) []*OrderBook {
	if pb == nil {
		return nil
	}
	ob := make([]*OrderBook, len(pb))
	for i, v := range pb {
		ob[i] = orderBookFromPB(v)
	}
	return ob
}

// 实时摆盘
type RTOrderBook struct {
	Security                *Security    //*股票
	Asks                    []*OrderBook //卖盘
	Bids                    []*OrderBook //买盘
	SvrRecvTimeBid          string       // 富途服务器从交易所收到数据的时间(for bid)部分数据的接收时间为零，例如服务器重启或第一次推送的缓存数据。该字段暂时只支持港股。
	SvrRecvTimeBidTimestamp float64      // 富途服务器从交易所收到数据的时间戳(for bid)
	SvrRecvTimeAsk          string       // 富途服务器从交易所收到数据的时间(for ask)
	SvrRecvTimeAskTimestamp float64      // 富途服务器从交易所收到数据的时间戳(for ask)
}

// K 线数据
type KLine struct {
	Time           string  //*时间戳字符串
	IsBlank        bool    //*是否是空内容的点,若为 true 则只有时间信息
	HighPrice      float64 //最高价
	OpenPrice      float64 //开盘价
	LowPrice       float64 //最低价
	ClosePrice     float64 //收盘价
	LastClosePrice float64 //昨收价
	Volume         int64   //成交量
	Turnover       float64 //成交额
	TurnoverRate   float64 //换手率（该字段为百分比字段，展示为小数表示）
	PE             float64 //市盈率
	ChangeRate     float64 //涨跌幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Timestamp      float64 //时间戳
}

func kLineFromPB(pb *qotcommon.KLine) *KLine {
	if pb == nil {
		return nil
	}
	return &KLine{
		Time:           pb.GetTime(),
		IsBlank:        pb.GetIsBlank(),
		HighPrice:      pb.GetHighPrice(),
		OpenPrice:      pb.GetOpenPrice(),
		LowPrice:       pb.GetLowPrice(),
		ClosePrice:     pb.GetClosePrice(),
		LastClosePrice: pb.GetLastClosePrice(),
		Volume:         pb.GetVolume(),
		Turnover:       pb.GetTurnover(),
		TurnoverRate:   pb.GetTurnoverRate(),
		PE:             pb.GetPe(),
		ChangeRate:     pb.GetChangeRate(),
		Timestamp:      pb.GetTimestamp(),
	}
}

func kLineListFromPB(pb []*qotcommon.KLine) []*KLine {
	if pb == nil {
		return nil
	}
	k := make([]*KLine, len(pb))
	for i, v := range pb {
		k[i] = kLineFromPB(v)
	}
	return k
}

// 分时数据
type TimeShare struct {
	Time           string  //*时间字符串
	Minute         int32   //*距离0点过了多少分钟
	IsBlank        bool    //*是否是空内容的点,若为 true 则只有时间信息
	Price          float64 //当前价
	LastClosePrice float64 //昨收价
	AvgPrice       float64 //均价
	Volume         int64   //成交量
	Turnover       float64 //成交额
	Timestamp      float64 //时间戳
}

func timeShareFromPB(pb *qotcommon.TimeShare) *TimeShare {
	if pb == nil {
		return nil
	}
	return &TimeShare{
		Time:           pb.GetTime(),
		Minute:         pb.GetMinute(),
		IsBlank:        pb.GetIsBlank(),
		Price:          pb.GetPrice(),
		LastClosePrice: pb.GetLastClosePrice(),
		AvgPrice:       pb.GetAvgPrice(),
		Volume:         pb.GetVolume(),
		Turnover:       pb.GetTurnover(),
		Timestamp:      pb.GetTimestamp(),
	}
}

func timeShareListFromPB(pb []*qotcommon.TimeShare) []*TimeShare {
	if pb == nil {
		return nil
	}
	t := make([]*TimeShare, len(pb))
	for i, v := range pb {
		t[i] = timeShareFromPB(v)
	}
	return t
}

// 实时分时
type RTData struct {
	Security   *Security    //*股票
	TimeShares []*TimeShare //*分时数据结构体
}

// 逐笔成交
type Ticker struct {
	Time         string                    //*时间字符串
	Sequence     int64                     //*唯一标识
	Dir          qotcommon.TickerDirection //*TickerDirection, 买卖方向
	Price        float64                   //*价格
	Volume       int64                     //*成交量
	Turnover     float64                   //*成交额
	RecvTime     float64                   //收到推送数据的本地时间戳，用于定位延迟
	Type         qotcommon.TickerType      //TickerType, 逐笔类型
	TypeSign     int32                     //逐笔类型符号
	PushDataType qotcommon.PushDataType    //用于区分推送情况，仅推送时有该字段
	Timestamp    float64                   //时间戳
}

func tickerFromPB(pb *qotcommon.Ticker) *Ticker {
	if pb == nil {
		return nil
	}
	return &Ticker{
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

func tickerListFromPB(pb []*qotcommon.Ticker) []*Ticker {
	if pb == nil {
		return nil
	}
	t := make([]*Ticker, len(pb))
	for i, v := range pb {
		t[i] = tickerFromPB(v)
	}
	return t
}

// 实时逐笔
type RTTicker struct {
	Security *Security
	Tickers  []*Ticker
}

// 买卖经纪
type Broker struct {
	ID   int64  //*经纪 ID
	Name string //*经纪名称
	Pos  int32  //*经纪档位
	//以下为 SF 行情特有字段
	OrderID int64 //交易所订单 ID，与交易接口返回的订单 ID 并不一样
	Volume  int64 //订单股数
}

func brokerFromPB(pb *qotcommon.Broker) *Broker {
	if pb == nil {
		return nil
	}
	return &Broker{
		ID:      pb.GetId(),
		Name:    pb.GetName(),
		Pos:     pb.GetPos(),
		OrderID: pb.GetOrderID(),
		Volume:  pb.GetVolume(),
	}
}

func brokerListFromPB(pb []*qotcommon.Broker) []*Broker {
	if pb == nil {
		return nil
	}
	b := make([]*Broker, len(pb))
	for i, v := range pb {
		b[i] = brokerFromPB(v)
	}
	return b
}

// 实时经纪队列
type BrokerQueue struct {
	Security *Security //*股票
	Asks     []*Broker //经纪 Ask(卖)盘
	Bids     []*Broker //经纪 Bid(买)盘
}
