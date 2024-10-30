package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/sirupsen/logrus"

	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
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

type config struct {
	certFile string
	keyFile  string
}

func initEnv() *config {
	cfg := &config{}

	if certFile := os.Getenv("TLS_CERT_FILE"); certFile != "" {
		cfg.certFile = certFile
	}

	if keyFile := os.Getenv("TLS_KEY_FILE"); keyFile != "" {
		cfg.keyFile = keyFile
	}

	return cfg
}

func initConfig() *ImageSwapConfig {
	// default location will be /etc/imageswap-config/imageswap.yaml
	imageSwapFilePath := "/etc/imageswap-config/imageswap.yaml"
	cfg := &ImageSwapConfig{}

	content, err := os.ReadFile(imageSwapFilePath)
	if err != nil {
	}

	err = yaml.Unmarshal(content, cfg)

	return cfg
}

func main() {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	logger := kwhlogrus.NewLogrus(logrusLogEntry)

	cfg := initEnv()

	imageSwapCfg := initConfig()

	fmt.Println(imageSwapCfg)

	// Create our mutator
	mt := kwhmutating.MutatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		return swapPodMutator(imageSwapCfg, ctx, ar, obj)
	})

	mcfg := kwhmutating.WebhookConfig{
		ID:      "podAnnotate",
		Obj:     &corev1.Pod{},
		Mutator: mt,
		Logger:  logger,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook: %s", err)
		os.Exit(1)
	}

	// Get the handler for our webhook.
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: logger})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook handler: %s", err)
	}
	logger.Infof("Listening on :8080")
	err = http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, whHandler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serving webhook: %s", err)
	}

}
