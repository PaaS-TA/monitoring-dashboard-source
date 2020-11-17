package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMonitoring(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 1")
	RunSpecs(t, "Monitoring Suite")
}
