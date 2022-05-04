package controllers

import (
	"context"
	"reflect"
	"testing"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func TestJupyterLabReconciler_deploymentForJupyterLab(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx        context.Context
		jupyterlab *toolsv1alpha1.JupyterLab
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsv1.Deployment
	}{
		{
			name: "should create a deployment for a JupyterLab",
			fields: fields{
				Client: k8sClient,
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx: context.TODO(),
				jupyterlab: &toolsv1alpha1.JupyterLab{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-jupyterlab",
						Namespace: "default",
					},
					Spec: toolsv1alpha1.JupyterLabSpec{
						Image:   "jupyterlab/minimal",
						Version: "latest",
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-jupyterlab",
					Namespace: "default",
					Labels: map[string]string{
						"app":                          "jupyter",
						"chart":                        "test-jupyterlab",
						"app.kubernetes.io/managed-by": "jupyterlab-operator",
						"app.kubernetes.io/component":  "jupyterlab",
						"app.kubernetes.io/part-of":    "jupyterlab",
						"app.kubernetes.io/version":    "latest",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "jupyterlab"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": "jupyterlab"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "jupyterlab",
									Image: "jupyterlab/minimal:latest",
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											Protocol:      corev1.ProtocolTCP,
											ContainerPort: 8888,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JupyterLabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.deploymentForJupyterLab(tt.args.ctx, tt.args.jupyterlab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JupyterLabReconciler.deploymentForJupyterLab() = %v, want %v", got, tt.want)
			}
		})
	}
}
