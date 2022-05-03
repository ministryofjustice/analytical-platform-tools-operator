package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

var _ = Describe("JupyterLab controller", func() {
	BeforeEach(func() {
		// failed test runs that don't clean up leave resources behind.
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		JupyterLabName      = "test-tool"
		JupyterLabNamespace = "default"
	)

	Context("create an instance of JupyterLab", func() {
		It("when a JupyterLab resource is submitted", func() {
			// Create a tool object
			tool := &toolsv1alpha1.JupyterLab{
				ObjectMeta: metav1.ObjectMeta{
					Name:      JupyterLabName,
					Namespace: JupyterLabNamespace,
				},
				Spec: toolsv1alpha1.JupyterLabSpec{
					Image:   "jupyterlab/minimal:latest",
					Version: "latest",
				},
			}
			By("posting a tool to the api")
			Expect(k8sClient.Create(context.TODO(), tool)).Should(Succeed())
		})
	})
})
