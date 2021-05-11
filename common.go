package futuapi

import (
	"github.com/hurisheng/go-futu-api/pb/common"
)

type ProgramStatus struct {
	Type       common.ProgramStatusType //*当前状态
	StrExtDesc string                   //额外描述
}

func programStatusFromPB(pb *common.ProgramStatus) *ProgramStatus {
	if pb == nil {
		return nil
	}
	return &ProgramStatus{
		Type:       pb.GetType(),
		StrExtDesc: pb.GetStrExtDesc(),
	}
}
