// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestPrefixSuffixTransformer(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		PrepBuiltin("PrefixSuffixTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: notImportantHere
prefix: baked-
suffix: -pie
fieldSpecs:
  - path: metadata/name
`, `
apiVersion: v1
kind: Service
metadata:
  name: apple
spec:
  ports:
  - port: 7002
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crd
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: apiservice
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Service
metadata:
  name: baked-apple-pie
spec:
  ports:
  - port: 7002
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crd
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: apiservice
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: baked-cm-pie
`)

	rm = th.LoadAndRunTransformer(`
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: notImportantHere
prefix: test-
fieldSpecs:
  - kind: Deployment
    path: metadata/name
  - kind: Deployment
    path: spec/template/spec/containers/name
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
      - image: myapp
        name: main
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crd
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: apiservice
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  template:
    spec:
      containers:
      - image: myapp
        name: test-main
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crd
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: apiservice
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
`)

	rm = th.LoadAndRunTransformer(`
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: notImportantHere
prefix: test-
fieldSpecs:
  - kind: Secret
    path: metadata/name
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
      - image: myapp
        name: main
        envFrom:
        - secretRef:
            name: sec
---
apiVersion: v1
kind: Secret
metadata:
  name: sec
data:
  secretKey: c2VjcmV0VmFsdWUK
type: Opaque
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
      - envFrom:
        - secretRef:
            name: test-sec
        image: myapp
        name: main
---
apiVersion: v1
data:
  secretKey: c2VjcmV0VmFsdWUK
kind: Secret
metadata:
  name: test-sec
type: Opaque
`)
}
