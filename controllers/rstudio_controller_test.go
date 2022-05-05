package controllers

import (
	"context"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

var _ = Describe("Rstudio controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		RstudioName      = "rstudio-test"
		RstudioNamespace = "default"
	)

	Context("create an instance of Rstudio", func() {
		It("when a Rstudio resource is submitted", func() {
			// Create a tool object
			tool := &toolsv1alpha1.Rstudio{
				ObjectMeta: metav1.ObjectMeta{
					Name:      RstudioName,
					Namespace: RstudioNamespace,
				},
				Spec: toolsv1alpha1.RstudioSpec{
					Image:   "rocker/rstudio",
					Version: "latest",
				},
			}
			By("posting a tool to the api")
			Expect(k8sClient.Create(context.TODO(), tool)).Should(Succeed())
		})
	})
})

func TestRstudioReconciler_deploymentForRstudio(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx        context.Context
		jupyterlab *toolsv1alpha1.Rstudio
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsv1.Deployment
	}{
		{
			name: "should create a deployment for a Rstudio",
			fields: fields{
				Client: k8sClient,
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx: context.TODO(),
				jupyterlab: &toolsv1alpha1.Rstudio{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-rstudio",
						Namespace: "default",
					},
					Spec: toolsv1alpha1.RstudioSpec{
						Image:   "rocker/rstudio",
						Version: "latest",
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-rstudio",
					Namespace: "default",
					Labels: map[string]string{
						"app":                          "rstudio",
						"chart":                        "test-rstudio",
						"app.kubernetes.io/managed-by": "rstudio-operator",
						"app.kubernetes.io/component":  "rstudio",
						"app.kubernetes.io/part-of":    "rstudio",
						"app.kubernetes.io/version":    "latest",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "rstudio"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": "rstudio"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "rstudio",
									Image: "rocker/rstudio:latest",
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
			r := &RstudioReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.deploymentForRstudio(tt.args.ctx, tt.args.jupyterlab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RstudioReconciler.deploymentForRstudio() = %v, want %v", got, tt.want)
			}
		})
	}
}
