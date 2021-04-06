package futuapi

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hurisheng/go-futu-api/pb/common"
	"github.com/hurisheng/go-futu-api/pb/getglobalstate"
	"github.com/hurisheng/go-futu-api/pb/initconnect"
	"github.com/hurisheng/go-futu-api/pb/keepalive"
	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/protocol"
	"github.com/hurisheng/go-futu-api/tcp"
)

var (
	ErrInterrupted   = errors.New("process is interrupted")
	ErrChannelClosed = errors.New("channel is closed")
)

// FutuAPI 是富途开放API的主要操作对象。
type FutuAPI struct {
	// 连接配置，通过方法设置，不设置默认为零值
	clientVer int32
	clientID  string
	protoFmt  common.ProtoFmt
	connID    uint64

	// TCP连接，连接后设置
	conn *tcp.Conn
	// 数据接收注册表
	reg *protocol.Registry

	serial uint32
	mu     sync.Mutex
	// 发送心跳的定时器，连接后设置
	ticker *time.Ticker
	// 心跳定时器关闭信号通道
	done chan struct{}
}

// NewFutuAPI 创建API对象，并启动goroutine进行发送保活心跳.
func NewFutuAPI() *FutuAPI {
	return &FutuAPI{
		reg:  protocol.NewRegistry(),
		done: make(chan struct{}),
	}
}

// 设置调用接口信息, 非必调接口
func (api *FutuAPI) SetClientInfo(id string, ver int32) {
	api.clientID = id
	api.clientVer = ver
}

// 设置通讯协议 body 格式, 目前支持 Protobuf|Json 两种格式，默认 ProtoBuf, 非必调接口
func (api *FutuAPI) SetProtoFmt(fmt common.ProtoFmt) {
	api.protoFmt = fmt
}

// 获取连接 ID，连接初始化成功后才会有值
func (api *FutuAPI) ConnID() uint64 {
	return api.connID
}

// 连接FutuOpenD
func (api *FutuAPI) Connect(ctx context.Context, address string) error {
	conn, err := tcp.Dial("tcp", address, protocol.NewDecoder(api.reg))
	if err != nil {
		return err
	}
	api.conn = conn
	resp, err := api.initConnect(ctx, &initconnect.C2S{
		ClientVer:    &api.clientVer,
		ClientID:     &api.clientID,
		PushProtoFmt: (*int32)(&api.protoFmt),
	})
	if err != nil {
		return err
	}
	api.connID = resp.GetConnID()
	if d := resp.GetKeepAliveInterval(); d > 0 {
		api.ticker = time.NewTicker(time.Second * time.Duration(d))
		go api.heartBeat(ctx)
	}
	return nil
}

// 关闭连接
func (api *FutuAPI) Close(ctx context.Context) error {
	if err := api.conn.Close(); err != nil {
		return err
	}
	close(api.done)
	api.reg.Close()
	return nil
}

func (api *FutuAPI) heartBeat(ctx context.Context) {
	for {
		select {
		case <-api.done:
			api.ticker.Stop()
			return
		case <-api.ticker.C:
			if _, err := api.keepAlive(ctx, time.Now().Unix()); err != nil {
				return
			}
		}
	}
}

func (api *FutuAPI) send(proto uint32, req proto.Message, out protocol.RespChan) error {
	// 递增serial
	api.mu.Lock()
	defer api.mu.Unlock()
	// 在registry注册get channel
	if err := api.reg.AddGetChan(proto, api.serial+1, out); err != nil {
		return err
	}
	// 向服务器发送req
	if err := api.conn.Send(protocol.NewEncoder(proto, api.serial+1, req)); err != nil {
		if err := api.reg.RemoveChan(proto, api.serial+1); err != nil {
			return err
		}
		return err
	}
	api.serial++
	return nil
}

func (api *FutuAPI) subscribe(proto uint32, out protocol.RespChan) error {
	// 在registry注册update channel
	if err := api.reg.AddUpdateChan(proto, out); err != nil {
		return err
	}
	return nil
}

type ret interface {
	GetRetType() int32
	GetRetMsg() string
}

func result(r ret) error {
	if common.RetType(r.GetRetType()) != common.RetType_RetType_Succeed {
		return errors.New(r.GetRetMsg())
	}
	return nil
}

// InitConnect 初始化连接
func (api *FutuAPI) initConnect(ctx context.Context, req *initconnect.C2S) (*initconnect.S2C, error) {
	out := make(chan *initconnect.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return nil, err
	}
	if err := api.send(ProtoIDInitConnect, &initconnect.Request{C2S: req}, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-out:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), result(resp)
	}
}

// KeepAlive 保活心跳
func (api *FutuAPI) keepAlive(ctx context.Context, t int64) (int64, error) {
	out := make(chan *keepalive.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return 0, err
	}
	if err := api.send(ProtoIDKeepAlive, &keepalive.Request{C2S: &keepalive.C2S{
		Time: &t,
	}}, ch); err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		return 0, ErrInterrupted
	case resp, ok := <-out:
		if !ok {
			return 0, ErrChannelClosed
		}
		return resp.GetS2C().GetTime(), result(resp)
	}
}

type GlobalState struct {
	MarketHK       qotcommon.QotMarketState
	MarketUS       qotcommon.QotMarketState
	MarketSH       qotcommon.QotMarketState
	MarketSZ       qotcommon.QotMarketState
	MarketHKFuture qotcommon.QotMarketState
	MarketUSFuture qotcommon.QotMarketState
	QotLogined     bool
	TrdLogined     bool
	ServerVer      int32
	ServerBuildNo  int32
	Time           int64
	LocalTime      float64
	ProgramStatus  *ProgramStatus
	QotSvrIpAddr   string
	TrdSvrIpAddr   string
	ConnID         uint64
}

func globalStateFromPB(resp *getglobalstate.S2C) *GlobalState {
	if resp == nil {
		return nil
	}
	return &GlobalState{
		MarketHK:       qotcommon.QotMarketState(resp.GetMarketHK()),
		MarketUS:       qotcommon.QotMarketState(resp.GetMarketUS()),
		MarketSH:       qotcommon.QotMarketState(resp.GetMarketSH()),
		MarketSZ:       qotcommon.QotMarketState(resp.GetMarketSZ()),
		MarketHKFuture: qotcommon.QotMarketState(resp.GetMarketHKFuture()),
		MarketUSFuture: qotcommon.QotMarketState(resp.GetMarketUSFuture()),
		QotLogined:     resp.GetQotLogined(),
		TrdLogined:     resp.GetTrdLogined(),
		ServerVer:      resp.GetServerVer(),
		ServerBuildNo:  resp.GetServerBuildNo(),
		Time:           resp.GetTime(),
		LocalTime:      resp.GetLocalTime(),
		ProgramStatus:  programStatusFromPB(resp.GetProgramStatus()),
		QotSvrIpAddr:   resp.GetQotSvrIpAddr(),
		TrdSvrIpAddr:   resp.GetTrdSvrIpAddr(),
		ConnID:         resp.GetConnID(),
	}
}

// 获取全局状态
func (api *FutuAPI) GetGlobalState(ctx context.Context) (*GlobalState, error) {
	out := make(chan *getglobalstate.Response)
	ch, err := protocol.NewPBChan(out)
	if err != nil {
		return nil, err
	}
	var userID uint64
	if err := api.send(ProtoIDGetGlobalState, &getglobalstate.Request{C2S: &getglobalstate.C2S{
		UserID: &userID,
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
		return globalStateFromPB(resp.GetS2C()), result(resp)
	}
}

// 系统推送通知
// func (api *FutuAPI) Notify() (<-chan *notify.Response, error) {
// 	out := make(chan *notify.Response)
// 	ch, err := protocol.NewPBChan(out)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := api.subscribe(ProtoIDNotify, ch); err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }
