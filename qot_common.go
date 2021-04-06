package futuapi

import "github.com/hurisheng/go-futu-api/pb/qotcommon"

type Security struct {
	Market qotcommon.QotMarket
	Code   string
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

type SubInfo struct {
	SubType    qotcommon.SubType
	Securities []*Security
}

func subInfoFromPB(pb *qotcommon.SubInfo) *SubInfo {
	if pb == nil {
		return nil
	}
	info := SubInfo{
		SubType: qotcommon.SubType(pb.GetSubType()),
	}
	if list := pb.GetSecurityList(); list != nil {
		info.Securities = make([]*Security, len(list))
		for i, v := range list {
			info.Securities[i] = securityFromPB(v)
		}
	}
	return &info
}

type ConnSubInfo struct {
	SubInfos  []*SubInfo
	UsedQuota int32
	IsOwnData bool
}

func connSubInfoFromPB(pb *qotcommon.ConnSubInfo) *ConnSubInfo {
	if pb == nil {
		return nil
	}
	info := ConnSubInfo{
		UsedQuota: pb.GetUsedQuota(),
		IsOwnData: pb.GetIsOwnConnData(),
	}
	if list := pb.GetSubInfoList(); list != nil {
		info.SubInfos = make([]*SubInfo, len(list))
		for i, v := range list {
			info.SubInfos[i] = subInfoFromPB(v)
		}
	}
	return &info
}

type OptionBasicQotExData struct {
	StrikePrice          float64                   //行权价
	ContractSize         int32                     //每份合约数(整型数据)
	ContractSizeFloat    float64                   //每份合约数（浮点型数据）
	OpenInterest         int32                     //未平仓合约数
	ImpliedVolatility    float64                   //隐含波动率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Premium              float64                   //溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%，如 20 实际对应 20%）
	Delta                float64                   //希腊值 Delta
	Gamma                float64                   //希腊值 Gamma
	Vega                 float64                   //希腊值 Vega
	Theta                float64                   //希腊值 Theta
	Rho                  float64                   //希腊值 Rho
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

type FutureBasicQotExData struct {
	LastSettlePrice    float64 //昨结
	Position           int32   //持仓量
	PositionChange     int32   //日增仓
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
