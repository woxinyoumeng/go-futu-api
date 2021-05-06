package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotrequesthistorykl"
)

const (
	ProtoIDQotRequestHistoryKL = 3103 //Qot_RequestHistoryKL	在线获取单只股票一段历史 K 线
)

// 获取历史 K 线
func (api *FutuAPI) RequestHistoryKLine(ctx context.Context, security *Security, begin string, end string, klType qotcommon.KLType, rehabType qotcommon.RehabType,
	maxNum int32, fields qotcommon.KLFields, nextKey []byte, extTime bool) (*HistoryKLine, error) {
	var klFields int64 = int64(fields)
	ch := make(qotrequesthistorykl.ResponseChan)
	if err := api.get(ProtoIDQotRequestHistoryKL, &qotrequesthistorykl.Request{
		C2S: &qotrequesthistorykl.C2S{
			RehabType:        (*int32)(&rehabType),
			KlType:           (*int32)(&klType),
			Security:         security.pb(),
			BeginTime:        &begin,
			EndTime:          &end,
			MaxAckKLNum:      &maxNum,
			NeedKLFieldsFlag: &klFields,
			NextReqKey:       nextKey,
			ExtendedTime:     &extTime,
		},
	}, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return historyKLineFromPB(resp.GetS2C()), result(resp)
	}
}

type HistoryKLine struct {
	Security *Security //证券
	KLines   []*KLine  //K 线数据
	NextKey  []byte    //分页请求 key。一次请求没有返回所有数据时，下次请求带上这个 key，会接着请求
}

func historyKLineFromPB(pb *qotrequesthistorykl.S2C) *HistoryKLine {
	if pb == nil {
		return nil
	}
	h := HistoryKLine{
		Security: securityFromPB(pb.GetSecurity()),
		KLines:   kLineListFromPB(pb.GetKlList()),
	}
	if list := pb.GetNextReqKey(); list != nil {
		h.NextKey = make([]byte, len(list))
		copy(h.NextKey, list)
	}
	return &h
}
