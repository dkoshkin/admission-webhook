package admitter

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodHostNetworkValidator checks pods
// Returns Allowed=false when hostNetwork=true
// Ignores any pods in 'kube-system' namespace
type PodHostNetworkValidator struct {
	UniversalDeserializerAdmitter
}

func (a PodHostNetworkValidator) Admit(ar *v1beta1.AdmissionReview) (*v1beta1.AdmissionResponse, error) {
	glog.V(2).Info("admitting pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		err := fmt.Errorf("expected resource to be %s", podResource)
		return nil, err
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := a.Decoder.Decode(raw, nil, &pod); err != nil {
		err := fmt.Errorf("could not decode pod: %v", err)
		return nil, err
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	// ignore any pods in kube-system
	if pod.GetObjectMeta().GetNamespace() == "kube-system" {
		return &reviewResponse, nil
	}

	// check pod
	if pod.Spec.HostNetwork {
		reviewResponse.Allowed = false
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace("the pod is deployed in hostNetwork")}
	}

	return &reviewResponse, nil
}
