/*
Copyright 2022.

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

package classification_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/projectsveltos/classifier-agent/pkg/classification"
	"github.com/projectsveltos/classifier-agent/pkg/utils"
	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
	libsveltosutils "github.com/projectsveltos/libsveltos/lib/utils"
)

const (
	version25 = "v1.25.2"
	version24 = "v1.24.2"
)

var (
	podTemplate = `apiVersion: v1
kind: Pod
metadata:
  namespace: %s
  name: %s
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2`
)

var _ = Describe("Manager: evaluation", func() {
	var scheme *runtime.Scheme

	BeforeEach(func() {
		var err error
		scheme, err = setupScheme()
		Expect(err).ToNot(HaveOccurred())
		classification.Reset()
	})

	It("IsVersionAMatch: comparisonEqual returns true when version matches", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonEqual)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: comparisonEqual returns false when version does not match", func() {
		node := getNode(version24)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonEqual)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: ComparisonNotEqual returns true when version doesn't match", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonNotEqual)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: ComparisonNotEqual returns false when version matches", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonNotEqual)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: ComparisonGreaterThan returns true when version is strictly greater than specified one", func() {
		node := getNode("v1.25.4")
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonGreaterThan)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: ComparisonGreaterThan returns false when version is not strictly greater than specified one", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonGreaterThan)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: ComparisonGreaterThanOrEqualTo returns true when version is equal to specified one", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonGreaterThanOrEqualTo)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: ComparisonGreaterThanOrEqualTo returns false when version is not equal/greater than specified one", func() {
		node := getNode(version24)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonGreaterThanOrEqualTo)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: ComparisonLessThan returns true when version is strictly less than specified one", func() {
		node := getNode("v1.25.0")
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonLessThan)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: ComparisonLessThan returns false when version is not strictly less than specified one", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonLessThan)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: ComparisonLessThanOrEqualTo returns true when version is equal to specified one", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version25, libsveltosv1alpha1.ComparisonLessThanOrEqualTo)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: ComparisonLessThanOrEqualTo returns false when version is not equal/less than specified one", func() {
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonLessThanOrEqualTo)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("IsVersionAMatch: multiple constraints returns true when all checks pass", func() {
		node := getNode(version24)
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThanOrEqualTo)
		classifier.Spec.KubernetesVersionConstraints = append(classifier.Spec.KubernetesVersionConstraints,
			libsveltosv1alpha1.KubernetesVersionConstraint{
				Comparison: string(libsveltosv1alpha1.ComparisonLessThan),
				Version:    version25,
			},
		)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeTrue())
	})

	It("IsVersionAMatch: multiple constraints returns false if at least one check fails", func() {
		node := getNode(version24)
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThan)
		classifier.Spec.KubernetesVersionConstraints = append(classifier.Spec.KubernetesVersionConstraints,
			libsveltosv1alpha1.KubernetesVersionConstraint{
				Comparison: string(libsveltosv1alpha1.ComparisonLessThan),
				Version:    version25,
			},
		)

		initObjects := []client.Object{
			node,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		match, err := classification.IsVersionAMatch(manager, context.TODO(),
			classifier)
		Expect(err).To(BeNil())
		Expect(match).To(BeFalse())
	})

	It("createClassifierReport creates ClassifierReport", func() {
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThan)

		initObjects := []client.Object{
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		isMatch := true
		Expect(classification.CreateClassifierReport(manager, context.TODO(), classifier, isMatch)).To(Succeed())

		verifyClassifierReport(c, classifier, isMatch)
	})

	It("createClassifierReport updates ClassifierReport", func() {
		phase := libsveltosv1alpha1.ReportProcessed
		isMatch := false
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThan)
		classifierReport := &libsveltosv1alpha1.ClassifierReport{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: utils.ReportNamespace,
				Name:      classifier.Name,
			},
			Spec: libsveltosv1alpha1.ClassifierReportSpec{
				ClassifierName: classifier.Name,
				Match:          !isMatch,
			},
			Status: libsveltosv1alpha1.ClassifierReportStatus{
				Phase: &phase,
			},
		}
		initObjects := []client.Object{
			classifierReport,
			classifier,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		Expect(classification.CreateClassifierReport(manager, context.TODO(), classifier, isMatch)).To(Succeed())

		verifyClassifierReport(c, classifier, isMatch)
	})

	It("evaluateClassifierInstance creates ClassifierReport", func() {
		// Create node and classifier so cluster is a match
		isMatch := true
		node := getNode(version25)
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThan)

		initObjects := []client.Object{
			classifier,
			node,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		Expect(classification.EvaluateClassifierInstance(manager, context.TODO(), classifier.Name)).To(Succeed())

		verifyClassifierReport(c, classifier, isMatch)
	})

	It("isResourceAMatch returns true when resources are match for classifier", func() {
		countMin := 3
		countMax := 5
		namespace := randomString()
		classifier := &libsveltosv1alpha1.Classifier{
			ObjectMeta: metav1.ObjectMeta{
				Name: randomString(),
			},
			Spec: libsveltosv1alpha1.ClassifierSpec{
				ClassifierLabels: []libsveltosv1alpha1.ClassifierLabel{
					{Key: randomString(), Value: randomString()},
				},
				DeployedResourceConstraints: []libsveltosv1alpha1.DeployedResourceConstraint{
					{
						Namespace: namespace,
						MinCount:  &countMin,
						MaxCount:  &countMax,
						Group:     "",
						Version:   "v1",
						Kind:      "Pod",
					},
				},
			},
		}

		Expect(testEnv.Create(context.TODO(), classifier)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, classifier)).To(Succeed())

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		Expect(testEnv.Create(context.TODO(), ns)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, ns)).To(Succeed())

		for i := 0; i < countMin-1; i++ {
			pod := fmt.Sprintf(podTemplate, namespace, randomString())
			u, err := libsveltosutils.GetUnstructured([]byte(pod))
			Expect(err).To(BeNil())
			Expect(testEnv.Create(context.TODO(), u)).To(Succeed())
			Expect(waitForObject(context.TODO(), testEnv.Client, u)).To(Succeed())
		}

		watcherCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		classification.InitializeManager(watcherCtx, klogr.New(), testEnv.Config, testEnv.Client, nil, 10)
		manager := classification.GetManager()

		isMatch, err := classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
		Expect(err).To(BeNil())
		Expect(isMatch).To(BeFalse())

		// Add one more pod so now the number of pod matches countMin
		pod := fmt.Sprintf(podTemplate, namespace, randomString())
		u, err := libsveltosutils.GetUnstructured([]byte(pod))
		Expect(err).To(BeNil())
		Expect(testEnv.Create(context.TODO(), u)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, u)).To(Succeed())

		isMatch, err = classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
		Expect(err).To(BeNil())
		Expect(isMatch).To(BeTrue())

		// Add more pods so there are now more than maxcount
		for i := 0; i < countMax; i++ {
			pod := fmt.Sprintf(podTemplate, namespace, randomString())
			var u *unstructured.Unstructured
			u, err = libsveltosutils.GetUnstructured([]byte(pod))
			Expect(err).To(BeNil())
			Expect(testEnv.Create(context.TODO(), u)).To(Succeed())
			Expect(waitForObject(context.TODO(), testEnv.Client, u)).To(Succeed())
		}

		isMatch, err = classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
		Expect(err).To(BeNil())
		Expect(isMatch).To(BeFalse())
	})

	It("isResourceAMatch returns true when resources with labels are match for classifier", func() {
		countMin := 1
		namespace := randomString()
		key1 := randomString()
		value1 := randomString()
		key2 := randomString()
		value2 := randomString()
		classifier := &libsveltosv1alpha1.Classifier{
			ObjectMeta: metav1.ObjectMeta{
				Name: randomString(),
			},
			Spec: libsveltosv1alpha1.ClassifierSpec{
				ClassifierLabels: []libsveltosv1alpha1.ClassifierLabel{
					{Key: randomString(), Value: randomString()},
				},
				DeployedResourceConstraints: []libsveltosv1alpha1.DeployedResourceConstraint{
					{
						Namespace: namespace,
						LabelFilters: []libsveltosv1alpha1.LabelFilter{
							{Key: key1, Operation: libsveltosv1alpha1.OperationEqual, Value: value1},
							{Key: key2, Operation: libsveltosv1alpha1.OperationEqual, Value: value2},
						},
						MinCount: &countMin,
						Group:    "",
						Version:  "v1",
						Kind:     "Pod",
					},
				},
			},
		}

		Expect(testEnv.Create(context.TODO(), classifier)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, classifier)).To(Succeed())

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		Expect(testEnv.Create(context.TODO(), ns)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, ns)).To(Succeed())

		// Create enough pod to match CountMin. But labels are currently not a match for classifier
		for i := 0; i <= countMin; i++ {
			pod := fmt.Sprintf(podTemplate, namespace, randomString())
			u, err := libsveltosutils.GetUnstructured([]byte(pod))
			u.SetLabels(map[string]string{key1: value1})
			Expect(err).To(BeNil())
			Expect(testEnv.Create(context.TODO(), u)).To(Succeed())
			Expect(waitForObject(context.TODO(), testEnv.Client, u)).To(Succeed())
		}

		watcherCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		classification.InitializeManager(watcherCtx, klogr.New(), testEnv.Config, testEnv.Client, nil, 10)
		manager := classification.GetManager()

		isMatch, err := classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
		Expect(err).To(BeNil())
		Expect(isMatch).To(BeFalse())

		// Get all pods in namespace and add second label needed for pods to match classifier
		podList := &corev1.PodList{}
		listOptions := []client.ListOption{
			client.InNamespace(namespace),
		}
		Expect(testEnv.List(context.TODO(), podList, listOptions...)).To(Succeed())

		for i := range podList.Items {
			pod := &podList.Items[i]
			pod.Labels[key2] = value2
			Expect(testEnv.Update(context.TODO(), pod)).To(Succeed())
		}

		// Use Eventually so cache is in sync
		Eventually(func() bool {
			isMatch, err = classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
			return err == nil && isMatch
		}, timeout, pollingInterval).Should(BeTrue())
	})

	It("isResourceAMatch returns true when fields resources are match for classifier", func() {
		countMin := 1
		namespace := randomString()
		podIP := "192.168.10.1"
		classifier := &libsveltosv1alpha1.Classifier{
			ObjectMeta: metav1.ObjectMeta{
				Name: randomString(),
			},
			Spec: libsveltosv1alpha1.ClassifierSpec{
				ClassifierLabels: []libsveltosv1alpha1.ClassifierLabel{
					{Key: randomString(), Value: randomString()},
				},
				DeployedResourceConstraints: []libsveltosv1alpha1.DeployedResourceConstraint{
					{
						Namespace: namespace,
						FieldFilters: []libsveltosv1alpha1.FieldFilter{
							{Field: "status.podIP", Operation: libsveltosv1alpha1.OperationEqual, Value: podIP},
						},
						MinCount: &countMin,
						Group:    "",
						Version:  "v1",
						Kind:     "Pod",
					},
				},
			},
		}

		Expect(testEnv.Create(context.TODO(), classifier)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, classifier)).To(Succeed())

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		Expect(testEnv.Create(context.TODO(), ns)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, ns)).To(Succeed())

		// fields are not a match for classifier
		pod := fmt.Sprintf(podTemplate, namespace, randomString())
		u, err := libsveltosutils.GetUnstructured([]byte(pod))
		Expect(err).To(BeNil())
		Expect(testEnv.Create(context.TODO(), u)).To(Succeed())
		Expect(waitForObject(context.TODO(), testEnv.Client, u)).To(Succeed())

		watcherCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		classification.InitializeManager(watcherCtx, klogr.New(), testEnv.Config, testEnv.Client, nil, 10)
		manager := classification.GetManager()

		isMatch, err := classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
		Expect(err).To(BeNil())
		Expect(isMatch).To(BeFalse())

		// Get all pods in namespace and add second label needed for pods to match classifier
		podList := &corev1.PodList{}
		listOptions := []client.ListOption{
			client.InNamespace(namespace),
		}
		Expect(testEnv.List(context.TODO(), podList, listOptions...)).To(Succeed())

		for i := range podList.Items {
			pod := &podList.Items[i]
			pod.Status.PodIP = podIP
			Expect(testEnv.Status().Update(context.TODO(), pod)).To(Succeed())
		}

		// Use Eventually so cache is in sync
		Eventually(func() bool {
			isMatch, err = classification.IsResourceAMatch(manager, watcherCtx, &classifier.Spec.DeployedResourceConstraints[0])
			return err == nil && isMatch
		}, timeout, pollingInterval).Should(BeTrue())
	})

	It("cleanClassifierReport removes classifier", func() {
		classifier := getClassifierWithKubernetesConstraints(version24, libsveltosv1alpha1.ComparisonGreaterThan)
		classifierReport := &libsveltosv1alpha1.ClassifierReport{
			ObjectMeta: metav1.ObjectMeta{
				Name:      classifier.Name,
				Namespace: utils.ReportNamespace,
			},
		}

		initObjects := []client.Object{
			classifierReport,
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()

		classification.InitializeManagerWithSkip(context.TODO(), klogr.New(), nil, c, nil, 10)

		manager := classification.GetManager()
		Expect(manager).ToNot(BeNil())

		Expect(classification.CleanClassifierReport(manager, context.TODO(), classifier.Name)).To(Succeed())

		err := c.Get(context.TODO(), types.NamespacedName{Name: classifier.Name, Namespace: utils.ReportNamespace}, classifier)
		Expect(err).ToNot(BeNil())
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
	})
})

func verifyClassifierReport(c client.Client, classifier *libsveltosv1alpha1.Classifier, isMatch bool) {
	classifierReport := &libsveltosv1alpha1.ClassifierReport{}
	Expect(c.Get(context.TODO(), types.NamespacedName{Namespace: utils.ReportNamespace, Name: classifier.Name},
		classifierReport)).To(Succeed())
	Expect(classifierReport.Labels).ToNot(BeNil())
	v, ok := classifierReport.Labels[libsveltosv1alpha1.ClassifierLabelName]
	Expect(ok).To(BeTrue())
	Expect(v).To(Equal(classifier.Name))
	Expect(classifierReport.Spec.Match).To(Equal(isMatch))
	Expect(classifierReport.Spec.ClassifierName).To(Equal(classifier.Name))
	Expect(classifierReport.Status.Phase).ToNot(BeNil())
	Expect(*classifierReport.Status.Phase).To(Equal(libsveltosv1alpha1.ReportWaitingForDelivery))
}