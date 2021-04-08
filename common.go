package futuapi

import (
	"github.com/hurisheng/go-futu-api/pb/common"
)

type CommonProgramStatus struct {
	Type       common.ProgramStatusType //*当前状态
	StrExtDesc string                   //额外描述
}

func commonProgramStatusFromPB(pb *common.ProgramStatus) *CommonProgramStatus {
	if pb == nil {
		return nil
	}
	return &CommonProgramStatus{
		Type:       pb.GetType(),
		StrExtDesc: pb.GetStrExtDesc(),
	}
}
