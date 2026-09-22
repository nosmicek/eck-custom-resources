package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	k8sv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Event struct {
	Object  runtime.Object
	Name    string
	Reason  string
	Message string
}

type ErrorEvent struct {
	Event
	Err error
}

// SkipRemoteDeleteAnnotation opts an object out of deleting its counterpart in
// Elasticsearch or Kibana. The remote object is left in place and only the
// Kubernetes object goes away.
const SkipRemoteDeleteAnnotation = "eck.github.com/skip-remote-delete"

func SkipRemoteDelete(o client.Object) bool {
	return o.GetAnnotations()[SkipRemoteDeleteAnnotation] == "true"
}

// ReleaseFinalizer drops the finalizer without contacting the remote system, for the
// cases where that call can never succeed: the target instance is gone, the
// reconciler is disabled, or the object opted out. Without it the object would sit in
// Terminating forever, which also blocks deletion of its namespace.
func ReleaseFinalizer(ctx context.Context, cli client.Client, o client.Object, finalizer string) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		return ctrl.Result{}, nil
	}

	controllerutil.RemoveFinalizer(o, finalizer)
	if err := cli.Update(ctx, o); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func GetRequeueResult() ctrl.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(time.Duration.Minutes(1)),
	}
}

func RecordError(recorder record.EventRecorder, errorEvent ErrorEvent) {
	recorder.Event(errorEvent.Object, "Warning", errorEvent.Reason,
		fmt.Sprintf("%s for %s: %s", errorEvent.Message, errorEvent.Name, errorEvent.Err.Error()))
}

func RecordSuccess(recorder record.EventRecorder, event Event) {
	message := fmt.Sprintf("%s successful for %s", event.Reason, event.Name)
	if event.Message != "" {
		message = fmt.Sprintf("%s for %s", event.Message, event.Name)
	}

	recorder.Event(event.Object, "Normal", event.Reason, message)
}

func RecordEventAndReturn(res ctrl.Result, err error, recorder record.EventRecorder, event Event) (ctrl.Result, error) {

	if err != nil {
		RecordError(recorder, ErrorEvent{
			Event: event,
			Err:   err,
		})
	} else {
		RecordSuccess(recorder, event)
	}

	return res, err
}

// ErrSecretNotFound reports that a secret a resource points at is absent, as opposed
// to any other reason a lookup can fail. Callers match on this rather than on a bare
// NotFound, which would also swallow a missing anything-else fetched along the way.
var ErrSecretNotFound = errors.New("referenced secret not found")

func getSecret(cli client.Client, ctx context.Context, namespace string, name string, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("%w: %q in namespace %q", ErrSecretNotFound, name, namespace)
		}
		return err
	}
	return nil
}

func GetUserSecret(cli client.Client, ctx context.Context, namespace string, auth *configv2.UsernamePasswordAuthentication, secret *k8sv1.Secret) error {
	return getSecret(cli, ctx, namespace, auth.SecretName, secret)
}

func GetCertificateSecret(cli client.Client, ctx context.Context, namespace string, certificate *configv2.PublicCertificate, secret *k8sv1.Secret) error {
	return getSecret(cli, ctx, namespace, certificate.SecretName, secret)
}

func CommonEventFilter() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
	}
}
