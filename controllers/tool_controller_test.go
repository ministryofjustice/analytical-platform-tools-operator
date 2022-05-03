package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

var _ = Describe("Tools controller", func() {
	BeforeEach(func() {
		// failed test runs that don't clean up leave resources behind.
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		ToolName      = "test-tool"
		ToolNamespace = "default"
	)

	Context("create a tool", func() {
		It("when a tool resource is submitted", func() {
			// Create a tool object
			tool := &toolsv1alpha1.Tool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ToolName,
					Namespace: ToolNamespace,
				},
				Spec: toolsv1alpha1.ToolSpec{
					Image:        "test-image",
					Username:     "test-user",
					ImageVersion: "latest",
				},
			}
			By("posting a tool to the api")
			Expect(k8sClient.Create(context.TODO(), tool)).Should(Succeed())
			Expect(tool.Spec.Image).Should(Equal("test-image"))
			Expect(tool.Spec.Username).Should(Equal("test-user"))
			Expect(tool.Spec.ImageVersion).Should(Equal("latest"))
		})
	})
})
