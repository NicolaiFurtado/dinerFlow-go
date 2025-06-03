package models

import "google.golang.org/protobuf/types/known/timestamppb"

type StartShift struct {
	ID        int64                  `json:"id"`
	UserID    int64                  `json:"user_id"`
	StartTime *timestamppb.Timestamp `json:"start_time"`
	EndTime   *timestamppb.Timestamp `json:"end_time"` // nullable in DB
}
