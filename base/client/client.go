package client

import (
    "net/http"
)

type Client {
    addr string

}

func Dial(addr string) *Client {
    Client ret = Client{addr}

    return &ret
}

