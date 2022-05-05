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

var _ = Describe("airflow controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		airflowName      = "test-airflow"
		airflowNamespace = "default"
	)

	Context("create an instance of airflow", func() {
		It("when a airflow resource is submitted", func() {
			// Create a tool object
			tool := &toolsv1alpha1.Airflow{
				ObjectMeta: metav1.ObjectMeta{
					Name:      airflowName,
					Namespace: airflowNamespace,
				},
				Spec: toolsv1alpha1.AirflowSpec{
					Image:   "apache/airflow",
					Version: "latest",
				},
			}
			By("posting a tool to the api")
			Expect(k8sClient.Create(context.TODO(), tool)).Should(Succeed())
		})
	})
})

func TestAirflowReconciler_deploymentForairflow(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		airflow *toolsv1alpha1.Airflow
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsv1.Deployment
	}{
		{
			name: "should create a deployment for a airflow",
			fields: fields{
				Client: k8sClient,
				Scheme: scheme.Scheme,
			},
			args: args{
				airflow: &toolsv1alpha1.Airflow{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-airflow",
						Namespace: "default",
					},
					Spec: toolsv1alpha1.AirflowSpec{
						Image:   "apache/airflow",
						Version: "latest",
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-airflow",
					Namespace: "default",
					Labels: map[string]string{
						"app":                          "airflow",
						"chart":                        "test-airflow",
						"app.kubernetes.io/managed-by": "airflow-operator",
						"app.kubernetes.io/component":  "airflow",
						"app.kubernetes.io/part-of":    "airflow",
						"app.kubernetes.io/version":    "latest",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "airflow"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": "airflow"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "airflow",
									Image: "apache/airflow:latest",
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
			r := &AirflowReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.newDeploymentForAirflow(tt.args.airflow); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("airflowReconciler.deploymentForAirflow() = %v, want %v", got, tt.want)
			}
		})
	}
}
