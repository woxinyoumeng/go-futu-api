package futuapi

import (
	"github.com/hurisheng/go-futu-api/pb/common"
)

type ProgramStatus struct {
	Type       common.ProgramStatusType
	StrExtDesc string
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
