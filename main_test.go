package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestName(t *testing.T) {
	result, err := http.Get("http://localhost:8080/calculate?path=1-2-2&value=2-2-2")
	fmt.Printf("%v, %v\n", err, result)

	buf := make([]byte, 128)

	result.Body.Read(buf)
	resp := string(buf[0:result.ContentLength])
	fmt.Println(resp)
}
