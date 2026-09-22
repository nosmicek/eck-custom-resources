package utils

import (
	"context"
	"errors"
	"testing"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	k8sv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// A caller releases a finalizer on ErrSecretNotFound, so every other reason a lookup
// can fail must stay distinguishable from it - otherwise a transient or permission
// failure would be mistaken for "the secret is gone" and orphan the remote resource.
func TestGetUserSecretDistinguishesMissingFromOtherFailures(t *testing.T) {
	auth := &configv2.UsernamePasswordAuthentication{SecretName: "es-creds", UserName: "elastic"}

	t.Run("secret present", func(t *testing.T) {
		existing := &k8sv1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "es-creds", Namespace: "default"}}
		cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(existing).Build()

		var secret k8sv1.Secret
		if err := GetUserSecret(cli, context.Background(), "default", auth, &secret); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("secret missing is ErrSecretNotFound", func(t *testing.T) {
		cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		var secret k8sv1.Secret
		err := GetUserSecret(cli, context.Background(), "default", auth, &secret)
		if !errors.Is(err, ErrSecretNotFound) {
			t.Fatalf("expected ErrSecretNotFound, got %v", err)
		}
	})

	t.Run("forbidden is not ErrSecretNotFound", func(t *testing.T) {
		cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
			Get: func(_ context.Context, _ client.WithWatch, key client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
				return apierrors.NewForbidden(schema.GroupResource{Resource: "secrets"}, key.Name, errors.New("rbac denied"))
			},
		}).Build()

		var secret k8sv1.Secret
		err := GetUserSecret(cli, context.Background(), "default", auth, &secret)
		if err == nil {
			t.Fatal("expected an error")
		}
		if errors.Is(err, ErrSecretNotFound) {
			t.Fatalf("a forbidden lookup must not look like a missing secret, got %v", err)
		}
	})

	// A NotFound raised by something other than the secret lookup must not be
	// mistaken for a missing secret either - this is what a bare apierrors.IsNotFound
	// check at the call site would get wrong.
	t.Run("unrelated NotFound is not ErrSecretNotFound", func(t *testing.T) {
		unrelated := apierrors.NewNotFound(schema.GroupResource{Resource: "configmaps"}, "some-configmap")
		if errors.Is(unrelated, ErrSecretNotFound) {
			t.Fatal("an unrelated NotFound must not match ErrSecretNotFound")
		}
		if !apierrors.IsNotFound(unrelated) {
			t.Fatal("precondition: the unrelated error is a NotFound")
		}
	})
}

func TestGetCertificateSecretMissing(t *testing.T) {
	cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

	var secret k8sv1.Secret
	err := GetCertificateSecret(cli, context.Background(), "default",
		&configv2.PublicCertificate{SecretName: "es-ca", CertificateKey: "ca.crt"}, &secret)
	if !errors.Is(err, ErrSecretNotFound) {
		t.Fatalf("expected ErrSecretNotFound, got %v", err)
	}
}
