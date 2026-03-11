package test

import (
	"Smart_delivery_locker/utils/jwts"
	"fmt"
	"net/url"
	"testing"
)

func TestJwtP(test *testing.T) {
	claims, err := jwts.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IlJveSIsInJvbGUiOjEsInVzZXJJRCI6MSwiYXZhdGFyIjoiL3VwbG9hZHMvYXZhdGFyL-WktOWDjy5wbmciLCJleHAiOjE3NzMzNzIwNjMuODkwNjE4LCJpc3MiOiJMb2NrIn0.vkwpNkaribyBBLTyX1pnox6q4hYttKZTxjztpISQLmA")
	if err != nil {
		panic(err)
	}
	fmt.Println(claims)
}

func TestUrl(test *testing.T) {
	println(url.QueryEscape("李四"))
}
