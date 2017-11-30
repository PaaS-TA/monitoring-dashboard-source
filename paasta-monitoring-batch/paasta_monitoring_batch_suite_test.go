package paasta_monitoring_batch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPaastaMonitoringBatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PaastaMonitoringBatch Suite")
}
