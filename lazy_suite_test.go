package lazy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLazy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lazy Suite")
}
