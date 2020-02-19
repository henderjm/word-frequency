package main_test

import (
	"github.com/onsi/gomega/gexec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pathToWordCounter string

func TestWorkdayTest(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		pathToWordCounter, err = gexec.Build("github.com/henderjm/word-frequency")
		Expect(err).ToNot(HaveOccurred())
	})
	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "WordFrequency Suite")

}
