package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPaprika2markdown(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Paprika2markdown Suite")
}
