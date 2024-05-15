package main

import (
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
    "strings"
    "fmt"
)

func TestHandlerPost(t *testing.T){
    urls = make(map[string]string)
    
    type want struct {
        code int
      }
    tests := []struct {
        name string
        param string
        want want
    }{
        {
            name: "post test 1. body doesn't consist of data",
            param: "",
            want: want{
                code: 400,
            },
        },
        {
            name: "post test 2. body consist of data",
            param: "http://ya.ru",
            want: want{
                code: 201,
            },

        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            fmt.Printf("\n\nTest %v Body %v\n", test.name, test.param)
            param := strings.NewReader(test.param)
            request := httptest.NewRequest(http.MethodPost, "/", param)
            w := httptest.NewRecorder()
            handlerPost(w, request)

            res := w.Result()
			defer res.Body.Close()
            fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
            assert.Equal(t, test.want.code, res.StatusCode)
        })
    }
    
}

func TestHandlerGet(t *testing.T){
    type want struct {
        code int
    }
    tests := []struct {
        name string
        body string
        want want
    }{
        {
            name: "get test 1. body doesn't consist of data",
            body: "",
            want: want{
                code: 400,
            },
        },
        {
            name: "get test 2. body consist of data",
            body: "http://ya.ru",
            want: want{
                code: 307,
            },

        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            var addr string
            if test.body == "/"{
                addr =""
            } else {
                for k, v := range urls {
                    if v == test.body {
                        addr = k
                    }
                }
            }
            fmt.Printf("\n\nTest %v Body %v Addr %v\n", test.name, test.body, addr)
            request := httptest.NewRequest(http.MethodGet, "/"+addr, nil)
            w := httptest.NewRecorder()
            handlerGet(w, request)

            res := w.Result()
			defer res.Body.Close()
            fmt.Printf("want code = %d StatusCode %d\n", test.want.code, res.StatusCode)
            assert.Equal(t, test.want.code, res.StatusCode)
        })
    }

}