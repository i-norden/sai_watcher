package cup_actions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCupActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CupActions Suite")
}
