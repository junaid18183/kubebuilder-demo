/*
Copyright 2023 Juned Memon.

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
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	enbuildv1alpha1 "vivsoftorg/enbuild/api/v1alpha1"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
)

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//----------------------------------------------------------------------------------------------------------------------------------------------------------------
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/finalizers,verbs=update

//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=microservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=GitRepository,verbs=list;watch;get;patch;create;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)
	// TODO(user): your logic here

	application := &enbuildv1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "failed to get Cluster resource", "application", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// logger.Info("got application", "Owner", application.Spec.Owner)

	// Reconcile k8s gitrepository for application which holds the infra code
	// app_gitrepo := "https://github.com/" + application.Spec.Owner + "/" + application.Name + ".git"
	ReconcileGitRepositoryApplication(ctx, r, application, logger)

	// Reconcile microservices for application
	ReconcileMicroServiceApplication(ctx, r, application, logger)

	return ctrl.Result{}, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&enbuildv1alpha1.Application{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{}, predicate.AnnotationChangedPredicate{})).
		Owns(&sourcev1.GitRepository{}).
		Owns(&enbuildv1alpha1.MicroService{}).
		Complete(r)
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------

func ReconcileGitRepositoryApplication(ctx context.Context, r *ApplicationReconciler, application *enbuildv1alpha1.Application, logger logr.Logger) (ctrl.Result, error) {
	repo := CreateGitRepository(ctx, "application_"+application.Name, os.Getenv("GH_TOKEN"), true, application.Spec.Owner, false, "", "", logger)
	gitrepository, err := generateGitRepositorySpec(application.Name, application.Namespace, repo.GetHTMLURL(), application.Spec.SecretRef.Name, "main")

	if err != nil {
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(application, gitrepository, r.Scheme); err != nil {

		return ctrl.Result{}, err
	}
	if err := CreateGitRepositoryCR(ctx, r.Client, gitrepository, logger); err != nil {
		return ctrl.Result{}, err
	}

	application.Status.Repository = gitrepository.Spec.URL

	_err := r.Status().Update(ctx, application)
	if _err != nil {
		return ctrl.Result{}, _err
	}
	return ctrl.Result{}, _err
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------

func ReconcileMicroServiceApplication(ctx context.Context, r *ApplicationReconciler, application *enbuildv1alpha1.Application, logger logr.Logger) (ctrl.Result, error) {

	for i := range application.Spec.MicroServices {
		microservice, err := generateMicroServiceSpec(application.Spec.MicroServices[i].Name, application.Namespace, application.Spec.MicroServices[i].Spec.Template, application.Spec.SecretRef.Name, application.Spec.Owner)

		if err != nil {
			return ctrl.Result{}, err
		}

		if err := ctrl.SetControllerReference(application, microservice, r.Scheme); err != nil {

			return ctrl.Result{}, err
		}
		if err := CreateMicroService(ctx, r.Client, microservice, logger); err != nil {
			return ctrl.Result{}, err
		}

	}
	return ctrl.Result{}, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
