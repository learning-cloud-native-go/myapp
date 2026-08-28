package model

import (
	"fmt"
	"time"
)

type CivilDate time.Time

func (d CivilDate) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d CivilDate) MarshalJSON() ([]byte, error) {
	str := time.Time(d).Format("2006-01-02")
	return []byte(fmt.Sprintf("%q", str)), nil
}
