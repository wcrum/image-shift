/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
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
