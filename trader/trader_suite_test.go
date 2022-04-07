package trader_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trader Suite")
}
