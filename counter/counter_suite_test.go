package counter_test

import (
	"testing"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCounter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Counter Suite")
}

var _ = BeforeSuite(func() {
	httpmock.Activate()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})
