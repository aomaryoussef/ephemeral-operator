package controllers

import (
	"context"
	"fmt"
	"time"
// Kubernetes API imports
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"

	// Kubebuilder runtime and controller
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ephemeralv1 "github.com/aomaryoussef/ephemeral-operator/api/v1"
)

// EphemeralResourceReconciler reconciles an EphemeralResource object
type EphemeralResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is the main reconciliation loop
func (r *EphemeralResourceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// Fetch the EphemeralResource instance
	var ephemeralResource ephemeralv1.EphemeralResource
	if err := r.Get(ctx, req.NamespacedName, &ephemeralResource); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Check if TTL has exceeded the defined limit
	if time.Since(ephemeralResource.Status.LastModified.Time).Seconds() > float64(ephemeralResource.Spec.TTLSeconds) {
		// Delete resources that exceed TTL
		for _, resource := range ephemeralResource.Spec.Resources {
			if err := r.deleteResource(ctx, resource); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Update the LastModified field to track when this reconciliation happened
	ephemeralResource.Status.LastModified = metav1.NewTime(time.Now())
	if err := r.Status().Update(ctx, &ephemeralResource); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// deleteResource deletes a Kubernetes resource like Deployment, Secret, etc.
func (r *EphemeralResourceReconciler) deleteResource(ctx context.Context, resource ephemeralv1.ResourceRef) error {
	// Get the resource by its kind and name
	switch resource.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := r.Get(ctx, client.ObjectKey{Name: resource.Name}, &deployment); err != nil {
			return err
		}
		return r.Delete(ctx, &deployment)
	case "Secret":
		var secret corev1.Secret
		if err := r.Get(ctx, client.ObjectKey{Name: resource.Name}, &secret); err != nil {
			return err
		}
		return r.Delete(ctx, &secret)
	default:
		return fmt.Errorf("unsupported resource kind: %s", resource.Kind)
	}
}

// SetupWithManager sets up the controller with the Manager
func (r *EphemeralResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ephemeralv1.EphemeralResource{}).
		Owns(&corev1.Deployment{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
