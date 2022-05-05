/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

// AirflowReconciler reconciles a Airflow object
type AirflowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=airflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=airflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=airflows/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// +kubebuilder:rbac:groups=networking,resources=ingress,verbs=get;list;watch;create;update;patch;delete
// +

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *AirflowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the Airflow instance
	airflow := &toolsv1alpha1.Airflow{}
	err := r.Get(ctx, req.NamespacedName, airflow)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			log.Log.Info("Airflow resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	// Reconcile the Airflow instance
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, client.ObjectKey{Name: airflow.Name, Namespace: airflow.Namespace}, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("Airflow not found", "airflow", airflow.Name)
			dep := r.newDeploymentForAirflow(airflow)
			err = r.Create(ctx, dep)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		} else {
			log.Log.Error(err, "Failed to get Airflow")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *AirflowReconciler) newDeploymentForAirflow(airflow *toolsv1alpha1.Airflow) *appsv1.Deployment {
	image := airflow.Spec.Image
	version := airflow.Spec.Version

	if image == "" {
		image = "apache/airflow"
	}

	if version == "" {
		version = "latest"
	}

	lab := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      airflow.Name,
			Namespace: airflow.Namespace,
			Labels: map[string]string{
				"app":                          "airflow",
				"chart":                        airflow.Name,
				"app.kubernetes.io/managed-by": "airflow-operator",
				"app.kubernetes.io/component":  "airflow",
				"app.kubernetes.io/part-of":    "airflow",
				"app.kubernetes.io/version":    airflow.Spec.Version,
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
							Image: image + ":" + version,
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
	}
	ctrl.SetControllerReference(airflow, lab, r.Scheme)

	return lab
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirflowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Airflow{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
