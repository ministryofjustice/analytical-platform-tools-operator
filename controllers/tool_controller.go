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
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

// ToolReconciler reconciles a Tool object
type ToolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(jasonBirchall): Modify the Reconcile function to compare the state specified by
// the Tool object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ToolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Check the status of the tool
	tool := &toolsv1alpha1.Tool{}
	if err := r.Get(ctx, req.NamespacedName, tool); err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("tool resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Tool resource")
		return ctrl.Result{}, nil
	}

	if tool.Name == strings.ToLower("jupyterlab") {
		log.Log.Info("Reconciling JupyterLab")
		jupyterlab := &toolsv1alpha1.JupyterLab{}
		err := r.Get(ctx, types.NamespacedName{Name: tool.Name, Namespace: tool.Namespace}, jupyterlab)
		if err != nil && errors.IsNotFound(err) {
			jupyterlabDeploy := r.createJupyterLabDeployment(tool)
			err := r.Create(ctx, jupyterlabDeploy)
			if err != nil {
				log.Log.Error(err, "Failed to create JupyterLab resource")
				return ctrl.Result{}, nil
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}
	// TODO: Update status of the tool

	return ctrl.Result{}, nil
}

func (r *ToolReconciler) createJupyterLabDeployment(tool *toolsv1alpha1.Tool) *toolsv1alpha1.JupyterLab {
	jupyterlab := &toolsv1alpha1.JupyterLab{
		ObjectMeta: r.createObjectMeta(tool.Name, tool.Namespace),
		Spec: toolsv1alpha1.JupyterLabSpec{
			Image:   "jupyter/minimal-notebook",
			Version: tool.Spec.ImageVersion,
		},
	}

	ctrl.SetControllerReference(tool, jupyterlab, r.Scheme)
	return jupyterlab
}

func (r *ToolReconciler) createObjectMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ToolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Tool{}).
		Owns(&toolsv1alpha1.JupyterLab{}).
		Complete(r)
}
