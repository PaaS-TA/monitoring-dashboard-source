package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMonitoring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Monitoring Suite")
}