package main_test

import (
	"os/exec"

	"github.com/onsi/gomega/gbytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	Context("When missing arguments", func() {

		It("Should fail for missing required fields", func() {
			command := exec.Command(pathToWordCounter)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("the required flags `-i, --page-ids' and `-n, --number-of-words' were not specified"))
		})
	})

	Context("When n is 0 or lower", func() {

		It("Should error if n is zero", func() {
			command := exec.Command(pathToWordCounter, "-i", "some-id", "-n", "0")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("n must be greater than 0"))
		})

		It("Should error if n is negative", func() {
			command := exec.Command(pathToWordCounter, "-i", "some-id", "-n", "-1")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("n must be greater than 0"))
		})
	})
})
