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
)
