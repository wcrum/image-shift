package main

import (
	"context"

	"github.com/sirupsen/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func swapPodMutator(cfg *ImageSwapConfig, _ context.Context, _ *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// If not a pod just continue the mutation chain(if there is one) and don't do nothing.
		return &kwhmutating.MutatorResult{}, nil
	}

	// Mutate our object with the required annotations.

	// need to get current image

	// gcr.io::myregistry.com
	// credentials::--> credentials

	for i, container := range pod.Spec.Containers {
		img := cfg.ImageSwap.SwapImage(container.Image)

		if img != "" {
			pod.Spec.Containers[i].Image = img
			logrus.Infof("Patched %s with %s", container.Image, img)
		}
	}

	for i, container := range pod.Spec.InitContainers {
		img := cfg.ImageSwap.SwapImage(container.Image)

		if img != "" {
			pod.Spec.InitContainers[i].Image = img
			logrus.Infof("Patched %s with %s", container.Image, img)
		}

		logrus.Infof("%s --> %s", container.Image, img)
	}

	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	pod.Annotations["mutated"] = "true"
	pod.Annotations["mutator"] = "pod-annotate"

	return &kwhmutating.MutatorResult{
		MutatedObject: pod,
	}, nil
}
