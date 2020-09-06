package resty

import (
	"github.com/go-resty/resty/v2"
	"time"
)

var Client *resty.Client

func SetupResty() {
	Client = resty.New()
	Client.SetTimeout(30 * time.Second)
	Client.SetRetryCount(3)
}
