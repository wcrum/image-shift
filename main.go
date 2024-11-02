package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/sirupsen/logrus"

	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
)

func injectCertInMWC() {
	// creates the in-cluster config

	certFile := os.Getenv("TLS_CERT_FILE")

	fmt.Printf("Injecting Certificate!")

	crt, err := os.ReadFile(certFile)
	if err != nil {
		panic(err)
	}

	certEnc := base64.StdEncoding.EncodeToString([]byte(crt))
	fmt.Println(certEnc)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	mwc, err := clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.TODO(), "imageshift-webhook", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	for i := range mwc.Webhooks {
		fmt.Println(i)
		mwc.Webhooks[i].ClientConfig.CABundle = []byte(crt)
	}

	_, err = clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(context.TODO(), mwc, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Completed!")

}

func startWebhook() {
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

func verifyConnection() {
	// test connection to registry is stable
	// chose fastest repository maybe?

}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided.")
		return
	}

	fmt.Println(os.Args)

	// checking for app [init | webhook]
	// init mode injects certificate into the mutating webhook
	// webhook mode starts standarad webhook process

	// better way to do this? trying to minimalize dependcies.

	switch os.Args[1] {
	case "init":
		injectCertInMWC()
	case "webhook":
		startWebhook()
	case "testing":
		verifyConnection()
	}
}
