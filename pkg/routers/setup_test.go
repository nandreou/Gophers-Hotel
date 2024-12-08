package routers

import (
	"os"
	"testing"

	"nicksrepo.com/nick/pkg/config"
)

var app config.App

func SetUp(m *testing.M) {
	os.Exit(m.Run())
}
