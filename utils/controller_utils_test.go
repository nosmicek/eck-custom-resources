package utils

import (
	"context"
	"testing"

	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestSkipRemoteDelete(t *testing.T) {
	cases := []struct {
		name        string
		annotations map[string]string
		want        bool
	}{
		{name: "no annotations"},
		{name: "unrelated annotation", annotations: map[string]string{"other": "true"}},
		{name: "explicitly disabled", annotations: map[string]string{SkipRemoteDeleteAnnotation: "false"}},
		{name: "opted out of remote delete", annotations: map[string]string{SkipRemoteDeleteAnnotation: "true"}, want: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := &k8sv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Annotations: c.annotations}}
			if got := SkipRemoteDelete(o); got != c.want {
				t.Errorf("expected %v, got %v", c.want, got)
			}
		})
	}
}

func TestReleaseFinalizer(t *testing.T) {
	const finalizer = "test.eck.github.com/finalizer"
	const other = "other.eck.github.com/finalizer"

	t.Run("removes only its own finalizer", func(t *testing.T) {
		o := &k8sv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:       "cm",
			Namespace:  "default",
			Finalizers: []string{finalizer, other},
		}}
		cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(o).Build()

		if _, err := ReleaseFinalizer(context.Background(), cli, o, finalizer); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if controllerutil.ContainsFinalizer(o, finalizer) {
			t.Error("expected the finalizer to be removed")
		}
		if !controllerutil.ContainsFinalizer(o, other) {
			t.Error("expected an unrelated finalizer to be left in place")
		}
	})

	t.Run("no-op when the finalizer is absent", func(t *testing.T) {
		o := &k8sv1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:       "cm",
			Namespace:  "default",
			Finalizers: []string{other},
		}}
		cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(o).Build()

		if _, err := ReleaseFinalizer(context.Background(), cli, o, finalizer); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !controllerutil.ContainsFinalizer(o, other) {
			t.Error("expected an unrelated finalizer to be left in place")
		}
	})
}
