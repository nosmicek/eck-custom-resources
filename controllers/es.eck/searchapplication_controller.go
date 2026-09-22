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

package eseck

import (
	"context"
	"errors"
	"fmt"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"
	"k8s.io/client-go/tools/record"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// SearchApplicationReconciler reconciles a SearchApplication object
type SearchApplicationReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=searchapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=searchapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=searchapplications/finalizers,verbs=update

func (r *SearchApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "searchapplications.es.eck.github.com/finalizer"

	var searchApplication eseckv1alpha1.SearchApplication
	if err := r.Get(ctx, req.NamespacedName, &searchApplication); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	isDeleted := !searchApplication.ObjectMeta.DeletionTimestamp.IsZero()

	if isDeleted && utils.SkipRemoteDelete(&searchApplication) {
		logger.Info("Skipping remote deletion, releasing finalizer.", "Resource", req.NamespacedName)
		return utils.ReleaseFinalizer(ctx, r.Client, &searchApplication, finalizer)
	}

	targetInstance, err := r.getTargetInstance(&searchApplication, searchApplication.Spec.TargetConfig, ctx, req.Namespace)
	if err != nil {
		if isDeleted && apierrors.IsNotFound(err) {
			logger.Info("Target instance not found, releasing finalizer without remote deletion.", "Resource", req.NamespacedName)
			return utils.ReleaseFinalizer(ctx, r.Client, &searchApplication, finalizer)
		}
		return utils.GetRequeueResult(), err
	}

	if !targetInstance.Enabled {
		if isDeleted {
			logger.Info("Reconciler disabled, releasing finalizer without remote deletion.", "Resource", req.NamespacedName)
			return utils.ReleaseFinalizer(ctx, r.Client, &searchApplication, finalizer)
		}
		logger.Info("Elasticsearch reconciler disabled, not reconciling.", "Resource", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, *targetInstance, req)
	if createClientErr != nil {
		if isDeleted && errors.Is(createClientErr, utils.ErrSecretNotFound) {
			logger.Info("Referenced secret not found, releasing finalizer without remote deletion.", "Resource", req.NamespacedName, "Reason", createClientErr.Error())
			return utils.ReleaseFinalizer(ctx, r.Client, &searchApplication, finalizer)
		}
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	if searchApplication.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := esutils.DependenciesFulfilled(esClient, searchApplication.Spec.Dependencies); err != nil {
			r.Recorder.Event(&searchApplication, "Warning", "Missing dependencies",
				fmt.Sprintf("Some of declared dependencies are not present yet: %s", err.Error()))
			return utils.GetRequeueResult(), err
		}

		logger.Info("Creating/Updating search application", "searchApplication", req.Name)
		res, err := esutils.UpsertSearchApplication(esClient, searchApplication)

		if err == nil {
			r.Recorder.Event(&searchApplication, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", searchApplication.APIVersion, searchApplication.Kind, searchApplication.Name))
		} else {
			r.Recorder.Event(&searchApplication, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", searchApplication.APIVersion, searchApplication.Kind, searchApplication.Name, err.Error()))
		}

		if err := r.addFinalizer(&searchApplication, finalizer, ctx); err != nil {
			return ctrl.Result{}, err
		}
		return res, err
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&searchApplication, finalizer) {
			logger.Info("Deleting object", "searchApplication", searchApplication.Name)
			if _, err := esutils.DeleteSearchApplication(esClient, req.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&searchApplication, finalizer)
			if err := r.Update(ctx, &searchApplication); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SearchApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.SearchApplication{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}

func (r *SearchApplicationReconciler) addFinalizer(o client.Object, finalizer string, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

func (r *SearchApplicationReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
	targetInstance := r.ProjectConfig.Elasticsearch
	if TargetConfig.ElasticsearchInstance != "" {
		var resourceInstance eseckv1alpha1.ElasticsearchInstance
		if err := esutils.GetTargetElasticsearchInstance(r.Client, ctx, namespace, TargetConfig.ElasticsearchInstance, &resourceInstance); err != nil {
			r.Recorder.Event(object, "Warning", "Failed to load target instance", fmt.Sprintf("Target instance not found: %s", err.Error()))
			return nil, err
		}

		targetInstance = resourceInstance.Spec
	}
	return &targetInstance, nil
}
