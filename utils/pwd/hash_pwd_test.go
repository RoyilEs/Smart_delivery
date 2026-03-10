package pwd

import (
	"fmt"
	"testing"
)

var pwd = "1234"
var hashPwd string

func TestHashPassword(t *testing.T) {
	hashPwd = HashPassword(pwd)
	fmt.Println(hashPwd)
}

func TestComparePasswords(t *testing.T) {
	fmt.Println(ComparePasswords(hashPwd, pwd))
}
