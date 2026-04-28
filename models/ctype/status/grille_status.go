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
		str = "occupied"
	case GrilleDisabled:
		str = "grille_disabled"
	default:
		str = "nil"
	}
	return str
}
