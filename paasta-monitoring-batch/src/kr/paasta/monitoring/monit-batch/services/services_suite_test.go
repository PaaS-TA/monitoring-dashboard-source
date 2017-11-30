package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestServiceTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceTest Suite")
}
