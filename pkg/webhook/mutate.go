package webhook

import (
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"strings"
)

const (
	podsAnnotationPatch string = `[
		 {"op":"add","path":"/metadata/annotations","value":{"mutatedBy":"sample-webhook-server"}}
	]`
)

func mutatePods(ar v1.AdmissionReview) *v1.AdmissionResponse {
	shouldPatchPod := func(pod *corev1.Pod) bool {
		if value, found := pod.Labels["webhook-mutate"]; found {
			klog.V(2).Info("has mutate-me label with value:", value)
			return strings.ToLower(value) == "true"
		}
		return false
	}
	return applyPodPatch(ar, shouldPatchPod, podsAnnotationPatch)
}

func applyPodPatch(ar v1.AdmissionReview, shouldPatchPod func(*corev1.Pod) bool, patch string) *v1.AdmissionResponse {
	klog.V(2).Info("mutating pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		klog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if shouldPatchPod(&pod) {
		reviewResponse.Patch = []byte(patch)
		pt := v1.PatchTypeJSONPatch
		reviewResponse.PatchType = &pt
	}
	return &reviewResponse
}
