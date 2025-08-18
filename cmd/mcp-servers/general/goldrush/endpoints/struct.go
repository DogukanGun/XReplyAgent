package endpoints

import (
	"os"
)

type GoldrushEndpoints struct {
	BaseUrl   string
	AuthToken string
}

func NewGoldrushEndpoints() *GoldrushEndpoints {
	return &GoldrushEndpoints{
		BaseUrl:   "https://api.covalenthq.com/v1/",
		AuthToken: os.Getenv("GOLDRUSH_AUTH_TOKEN"),
	}
}
