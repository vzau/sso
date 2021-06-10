package main

import (
	"os"

	"github.com/microcosm-cc/bluemonday"
)

var policy *bluemonday.Policy

func Getenv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

func CreatePolicy() {
	policy = bluemonday.UGCPolicy()
}

func GetPolicy() *bluemonday.Policy {
	return policy
}
