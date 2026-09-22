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

// QueryRulesetReconciler reconciles a QueryRuleset object
type QueryRulesetReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=queryrulesets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=queryrulesets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=queryrulesets/finalizers,verbs=update

func (r *QueryRulesetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "queryrulesets.es.eck.github.com/finalizer"

	var queryRuleset eseckv1alpha1.QueryRuleset
	if err := r.Get(ctx, req.NamespacedName, &queryRuleset); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	isDeleted := !queryRuleset.ObjectMeta.DeletionTimestamp.IsZero()

	if isDeleted && utils.SkipRemoteDelete(&queryRuleset) {
		logger.Info("Skipping remote deletion, releasing finalizer.", "Resource", req.NamespacedName)
		return utils.ReleaseFinalizer(ctx, r.Client, &queryRuleset, finalizer)
	}

	targetInstance, err := r.getTargetInstance(&queryRuleset, queryRuleset.Spec.TargetConfig, ctx, req.Namespace)
	if err != nil {
		if isDeleted && apierrors.IsNotFound(err) {
			logger.Info("Target instance not found, releasing finalizer without remote deletion.", "Resource", req.NamespacedName)
			return utils.ReleaseFinalizer(ctx, r.Client, &queryRuleset, finalizer)
		}
		return utils.GetRequeueResult(), err
	}

	if !targetInstance.Enabled {
		if isDeleted {
			logger.Info("Reconciler disabled, releasing finalizer without remote deletion.", "Resource", req.NamespacedName)
			return utils.ReleaseFinalizer(ctx, r.Client, &queryRuleset, finalizer)
		}
		logger.Info("Elasticsearch reconciler disabled, not reconciling.", "Resource", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, *targetInstance, req)
	if createClientErr != nil {
		if isDeleted && errors.Is(createClientErr, utils.ErrSecretNotFound) {
			logger.Info("Referenced secret not found, releasing finalizer without remote deletion.", "Resource", req.NamespacedName, "Reason", createClientErr.Error())
			return utils.ReleaseFinalizer(ctx, r.Client, &queryRuleset, finalizer)
		}
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	if queryRuleset.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Creating/Updating query ruleset", "queryRuleset", req.Name)
		res, err := esutils.UpsertQueryRuleset(esClient, queryRuleset)

		if err == nil {
			r.Recorder.Event(&queryRuleset, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", queryRuleset.APIVersion, queryRuleset.Kind, queryRuleset.Name))
		} else {
			r.Recorder.Event(&queryRuleset, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", queryRuleset.APIVersion, queryRuleset.Kind, queryRuleset.Name, err.Error()))
		}

		if err := r.addFinalizer(&queryRuleset, finalizer, ctx); err != nil {
			return ctrl.Result{}, err
		}
		return res, err
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&queryRuleset, finalizer) {
			logger.Info("Deleting object", "queryRuleset", queryRuleset.Name)
			if _, err := esutils.DeleteQueryRuleset(esClient, req.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&queryRuleset, finalizer)
			if err := r.Update(ctx, &queryRuleset); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *QueryRulesetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.QueryRuleset{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}

func (r *QueryRulesetReconciler) addFinalizer(o client.Object, finalizer string, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

func (r *QueryRulesetReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
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
