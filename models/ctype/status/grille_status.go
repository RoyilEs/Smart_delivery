package status

import "github.com/goccy/go-json"

type GrilleStatus int

const (
	Idle           GrilleStatus = 0
	Occupied       GrilleStatus = 1
	GrilleDisabled GrilleStatus = -1
)

func (s GrilleStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s GrilleStatus) String() string {
	var str string
	switch s {
	case Idle:
		str = "idle"
	case Occupied:
		str = "occupied" // 占用
	case GrilleDisabled:
		str = "disabled" // 停用
	default:
		str = "nil"
	}
	return str
}
