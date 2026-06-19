package models

import (
	"fmt"
	"time"
)

type ErrorDTO struct {
	Msg  string
	Time time.Time
}

func (e *ErrorDTO) ToString() string {
	return fmt.Sprintf("[%s] error: %s", e.Time.Format(time.DateTime), e.Msg)
}
