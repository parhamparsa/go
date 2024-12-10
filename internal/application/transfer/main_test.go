package transfer_test

import (
	"app/internal/testbox"
	"testing"
)

var tb *testbox.TestBox //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	tb = testbox.SetupTestBox("test_app_transfer")
	m.Run()
}
