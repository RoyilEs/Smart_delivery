package status

import "github.com/goccy/go-json"

type UserStatus int

const (
	Enabled      UserStatus = 1
	UserDisabled UserStatus = 0
)

func (s UserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s UserStatus) String() string {
	var str string
	switch s {
	case Enabled:
		str = "enabled"
	case UserDisabled:
		str = "disabled"
	default:
		str = "nil"
	}
	return str
}
