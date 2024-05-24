/*
Copyright 2024 RamyaSurati.

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

package controller

import (
	"context"

	"client.Object"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	scanv1 "github.com/rsurati/mykubernetesproject/api/v1"
)

var obj client.Object

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=scan.core.mykubernetesproject.io,resources=clusterscans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scan.core.mykubernetesproject.io,resources=clusterscans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scan.core.mykubernetesproject.io,resources=clusterscans/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterScan object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = log.FromContext(ctx)

	// TODO(user): your logic here

	// Code
	// Fetch the ClusterScan instance
	clusterScan := &batchv1alpha1.ClusterScan{}
	err := r.Get(ctx, req.NamespacedName, clusterScan)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define the desired Job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterScan.Name + "-job",
			Namespace: clusterScan.Namespace,
		},
		Spec: clusterScan.Spec.JobTemplate,
	}

	// Set ClusterScan instance as the owner and controller
	if err := controllerutil.SetControllerReference(clusterScan, job, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the job already exists
	found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	if err != nil && client.IgnoreNotFound(err) == nil {
		log.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.Create(ctx, job)
		if err != nil {
			return ctrl.Result{}, err
		}
		clusterScan.Status.LastScheduleTime = &metav1.Time{Time: time.Now()}
		err = r.Status().Update(ctx, clusterScan)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scanv1.ClusterScan{}).
		Complete(r)
}
