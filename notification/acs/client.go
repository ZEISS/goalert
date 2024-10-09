package acs

import (
	"net"
	"net/http"
	"time"
)

var DefaultTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var DefaultClient = &http.Client{
	Transport: DefaultTransport,
	Timeout:   10 * time.Second,
}
