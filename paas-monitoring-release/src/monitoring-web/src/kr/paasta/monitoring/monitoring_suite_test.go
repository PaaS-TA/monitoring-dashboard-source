package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
)

func TestMonitoring(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 1")
	RunSpecs(t, "Monitoring Suite")
}

