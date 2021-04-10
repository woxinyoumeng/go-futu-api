package futuapi

//  for all the commands.
const (
	ProtoIDTrdGetAccList  = 2001 //Trd_GetAccList	获取业务账户列表
	ProtoIDTrdUnlockTrade = 2005 //Trd_UnlockTrade	解锁或锁定交易
	ProtoIDTrdSubAccPush  = 2008 //Trd_SubAccPush	订阅业务账户的交易推送数据

	ProtoIDTrdGetFunds        = 2101 //Trd_GetFunds	获取账户资金
	ProtoIDTrdGetPositionList = 2102 //Trd_GetPositionList	获取账户持仓
	ProtoIDTrdGetMaxTrdQtys   = 2111 //Trd_GetMaxTrdQtys	获取最大交易数量

	ProtoIDTrdGetOrderList            = 2201 //Trd_GetOrderList	获取订单列表
	ProtoIDTrdPlaceOrder              = 2202 //Trd_PlaceOrder	下单
	ProtoIDTrdModifyOrder             = 2205 //Trd_ModifyOrder
	ProtoIDTrdUpdateOrder             = 2208 //Trd_UpdateOrder	推送订单状态变动通知
	ProtoIDTrdGetOrderFillList        = 2211 //Trd_GetOrderFillList	获取成交列表
	ProtoIDTrdUpdateOrderFill         = 2218 //Trd_UpdateOrderFill	推送成交通知
	ProtoIDTrdGetHistoryOrderList     = 2221 //Trd_GetHistoryOrderList	获取历史订单列表
	ProtoIDTrdGetHistoryOrderFillList = 2222 //Trd_GetHistoryOrderFillList	获取历史成交列表

	ProtoIDQotUpdatePriceReminder = 3019 //Qot_UpdatePriceReminder	到价提醒通知

	ProtoIDQotRequestHistoryKL      = 3103 //Qot_RequestHistoryKL	在线获取单只股票一段历史 K 线
	ProtoIDQotRequestHistoryKLQuota = 3104 //Qot_RequestHistoryKLQuota	获取历史 K 线额度
	ProtoIDQotRequestRehab          = 3105 //Qot_RequestRehab	在线获取单只股票复权信息

	ProtoIDQotGetStaticInfo          = 3202 //Qot_GetStaticInfo	获取股票静态信息
	ProtoIDQotGetPlateSet            = 3204 //Qot_GetPlateSet	获取板块集合下的板块
	ProtoIDQotGetPlateSecurity       = 3205 //Qot_GetPlateSecurity	获取板块下的股票
	ProtoIDQotGetReference           = 3206 //Qot_GetReference	获取正股相关股票
	ProtoIDQotGetOwnerPlate          = 3207 //Qot_GetOwnerPlate	获取股票所属板块
	ProtoIDQotGetHoldingChangeList   = 3208 //Qot_GetHoldingChangeList	获取持股变化列表
	ProtoIDQotGetOptionChain         = 3209 //Qot_GetOptionChain	获取期权链
	ProtoIDQotGetWarrant             = 3210 //Qot_GetWarrant	获取窝轮
	ProtoIDQotGetCapitalFlow         = 3211 //Qot_GetCapitalFlow	获取资金流向
	ProtoIDQotGetCapitalDistribution = 3212 //Qot_GetCapitalDistribution
	ProtoIDQotGetUserSecurity        = 3213 //Qot_GetUserSecurity	获取自选股分组下的股票
	ProtoIDQotModifyUserSecurity     = 3214 //Qot_ModifyUserSecurity	修改自选股分组下的股票
	ProtoIDQotStockFilter            = 3215 //Qot_StockFilter	获取条件选股
	ProtoIDQotGetIpoList             = 3217 //Qot_GetIpoList	获取新股
	ProtoIDQotGetFutureInfo          = 3218 //Qot_GetFutureInfo	获取期货合约资料
	ProtoIDQotRequestTradeDate       = 3219 //Qot_RequestTradeDate	获取市场交易日，在线拉取不在本地计算
	ProtoIDQotSetPriceReminder       = 3220 //Qot_SetPriceReminder	设置到价提醒
	ProtoIDQotGetPriceReminder       = 3221 //Qot_GetPriceReminder	获取到价提醒
	ProtoIDQotGetUserSecurityGroup   = 3222 //Qot_GetUserSecurityGroup	获取自选股分组列表
	ProtoIDQotGetMarketState         = 3223 //Qot_GetMarketState	获取指定品种的市场状态
)
