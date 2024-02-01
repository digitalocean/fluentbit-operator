package plugins

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:object:generate:=true

type Secret struct {
	// ValueSource represents a source for the value of a secret.
	ValueFrom ValueSource `json:"valueFrom,omitempty"`
}

// +kubebuilder:object:generate:=true

type ValueSource struct {
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

type SecretLoader struct {
	client    client.Client
	namespace string
}

func NewSecretLoader(c client.Client, ns string) SecretLoader {
	return SecretLoader{
		client:    c,
		namespace: ns,
	}
}

func (sl SecretLoader) LoadSecret(s Secret) (string, error) {
	var secret corev1.Secret
	if err := sl.client.Get(context.Background(), client.ObjectKey{Name: s.ValueFrom.SecretKeyRef.Name, Namespace: sl.namespace}, &secret); err != nil {
		return "", err
	}

	if v, ok := secret.Data[s.ValueFrom.SecretKeyRef.Key]; !ok {
		return "", fmt.Errorf("The key %s is not found.", s.ValueFrom.SecretKeyRef.Key)
	} else {
		return strings.TrimSuffix(fmt.Sprintf("%s", v), "\n"), nil
	}
}
