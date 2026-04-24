package pwd

import (
	"fmt"
	"testing"
)

var pwd = "admin"
var hashPwd string = "$2a$04$QNvG84O62D2UbXJUmuhlE.HpajEHNp3w1iv.KEPvXT6c5HmDPGgvm"

func TestHashPassword(t *testing.T) {
	hashPwd = HashPassword(pwd)
	fmt.Println(hashPwd)
}

func TestComparePasswords(t *testing.T) {
	fmt.Println(ComparePasswords(hashPwd, pwd))
}
