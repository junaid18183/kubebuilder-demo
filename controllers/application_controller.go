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

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	enbuildv1alpha1 "vivsoftorg/enbuild/api/v1alpha1"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/finalizers,verbs=update

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

	// Reconcile k8s gitrepository.
	gitrepository, err := generateGitRepositorySpec(application, logger, r)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(application, gitrepository, r.Scheme); err != nil {

		return ctrl.Result{}, err
	}
	if err := CreateGitRepository(ctx, r, gitrepository, logger); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&enbuildv1alpha1.Application{}).
		Owns(&sourcev1.GitRepository{}).
		Complete(r)
}

func generateGitRepositorySpec(app *enbuildv1alpha1.Application, log logr.Logger, r *ApplicationReconciler) (*sourcev1.GitRepository, error) {

	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL:       "https://gitlab.com/enbuild-staging/iac-templates/bigbang",
			SecretRef: app.Spec.SecretRef,
			Reference: &sourcev1.GitRepositoryRef{Branch: "main"},
		},
	}, nil
}

func CreateGitRepository(ctx context.Context, r *ApplicationReconciler, gitrepository *sourcev1.GitRepository, log logr.Logger) error {
	foundGitRepository := &sourcev1.GitRepository{}
	if err := r.Get(ctx, types.NamespacedName{Name: gitrepository.Name, Namespace: gitrepository.Namespace}, foundGitRepository); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating GitRepository", "namespace", gitrepository.Namespace, "name", gitrepository.Name)
			if err := r.Create(ctx, gitrepository); err != nil {
				log.Error(err, "unable to create gitrepository")
				return err
			}
		} else {
			log.Error(err, "error getting gitrepository")
			return err
		}
	}

	return nil
}
