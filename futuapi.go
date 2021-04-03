package futuapi

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hurisheng/go-futu-api/protobuf/initconnect"
	"github.com/hurisheng/go-futu-api/protobuf/keepalive"
	"github.com/hurisheng/go-futu-api/protocol"
	"github.com/hurisheng/go-futu-api/tcp"
)

var (
	ErrInterrupted   = errors.New("process is interrupted")
	ErrChannelClosed = errors.New("channel is closed")
)

type Response interface {
	GetRetType() int32
	GetRetMsg() string
	GetErrCode() int32
}

// 获取方法返回的错误信息
type Failure struct {
	RetMsg  string //请求失败，说明失败原因
	ErrCode int32  //请求失败对应错误码
}

// 新建错误信息
func NewFailure(msg string, code int32) *Failure {
	return &Failure{
		RetMsg:  msg,
		ErrCode: code,
	}
}

// 实现error接口，返回失败原因
func (f *Failure) Error() string {
	return f.RetMsg
}

// FutuAPI 是富途开放API的主要操作对象。
type FutuAPI struct {
	conn *tcp.Conn
	reg  *protocol.Registry

	serial uint32
	mu     sync.Mutex
	ticker *time.Ticker
	done   chan struct{}
}

// Config 为API配置信息
type Config struct {
	Address   string
	ClientVer int32
	ClientID  string
}

// NewFutuAPI 创建API对象，并启动goroutine进行发送保活心跳.
func NewFutuAPI(config *Config) (*FutuAPI, error) {
	// connect socket
	reg := protocol.NewRegistry()
	conn, err := tcp.Dial("tcp", config.Address, protocol.NewDecoder(reg))
	if err != nil {
		return nil, fmt.Errorf("connect to server error: %w", err)
	}
	api := &FutuAPI{
		conn: conn,
		reg:  reg,
	}
	// init connect
	resp, err := api.initConnect(context.Background(), &initconnect.C2S{
		ClientVer: &config.ClientVer,
		ClientID:  &config.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("call InitConnect error: %w", err)
	}
	api.ticker = time.NewTicker(time.Second * time.Duration(resp.GetKeepAliveInterval()))
	go api.heartBeat()
	return api, nil
}

// Close 关闭连接.
func (api *FutuAPI) Close() error {
	if err := api.conn.Close(); err != nil {
		return err
	}
	close(api.done)
	api.reg.Close()
	return nil
}

func (api *FutuAPI) heartBeat() {
	for {
		select {
		case <-api.done:
			api.ticker.Stop()
			return
		case <-api.ticker.C:
			now := time.Now().Unix()
			_, err := api.keepAlive(context.Background(), &keepalive.C2S{
				Time: &now,
			},
			)
			if err != nil {
				return
			}
		}
	}
}

func (api *FutuAPI) send(proto uint32, req proto.Message, out protocol.ResponseChan) error {
	// 递增serial
	api.mu.Lock()
	api.serial++
	api.mu.Unlock()
	// 在registry注册get channel
	if err := api.reg.AddGetChan(proto, api.serial, out); err != nil {
		return err
	}
	// 向服务器发送req
	if err := api.conn.Send(protocol.NewEncoder(proto, api.serial, req)); err != nil {
		if err := api.reg.RemoveChan(proto, api.serial); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (api *FutuAPI) subscribe(proto uint32, out protocol.ResponseChan) error {
	// 在registry注册update channel
	if err := api.reg.AddUpdateChan(proto, out); err != nil {
		return err
	}
	return nil
}

// InitConnect 初始化连接
func (api *FutuAPI) initConnect(ctx context.Context, req *initconnect.C2S) (*initconnect.S2C, error) {
	ch := initconnect.NewResponseChan()
	if err := api.send(ProtoIDInitConnect, &initconnect.Request{C2S: req}, ch); err != nil {
		return nil, err
	}
	return ch.Response()
}

// KeepAlive 保活心跳
func (api *FutuAPI) keepAlive(ctx context.Context, req *keepalive.C2S) (*keepalive.S2C, error) {
	ch := keepalive.NewResponseChan()
	if err := api.send(ProtoIDKeepAlive, &keepalive.Request{C2S: req}, ch); err != nil {
		return nil, err
	}
	return ch.Response()
}

// // GetGlobalState 获取全局状态
// func (api *FutuAPI) GetGlobalState(req *GetGlobalState.Request) (<-chan *GetGlobalState.Response, error) {
// 	out := make(chan *GetGlobalState.Response)
// 	if err := api.send(ProtoIDGetGlobalState, req, out); err != nil {
// 		return nil, fmt.Errorf("GetGlobalState error: %w", err)
// 	}
// 	return out, nil
// }

// Notify 系统推送通知
// func (api *FutuAPI) Notify() (<-chan *Notify.Response, error) {
// 	out := make(Notify.ResponseChan)
// 	if err := api.subscribe(ProtoIDNotify, out); err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }
