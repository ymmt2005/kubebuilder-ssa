package controllers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ssav1 "github.com/ymmt2005/kubebuilder-ssa/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func myService() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(ssav1.GroupVersion.WithKind("MyService"))
	return obj
}

func dump(obj interface{}) {
	w := json.NewEncoder(os.Stderr)
	w.SetIndent("", "    ")
	w.Encode(obj)
}

var _ = Describe("Server Side Apply", func() {
	It("can create a CR", func() {
		obj := myService()
		obj.SetName("foo")
		obj.SetNamespace("default")
		obj.UnstructuredContent()["spec"] = map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{"protocol": "TCP", "port": 53},
				map[string]interface{}{"protocol": "UDP", "port": 53},
			},
		}

		err := k8sClient.Patch(context.TODO(), obj, client.Apply, &client.PatchOptions{
			FieldManager: "foo",
		})
		Expect(err).ShouldNot(HaveOccurred())

		foo := myService()
		err = k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "foo"}, foo)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can patch fields of the CR", func() {
		obj := myService()
		obj.SetName("foo")
		obj.SetNamespace("default")
		obj.UnstructuredContent()["spec"] = map[string]interface{}{
			"string":  "hoge",
			"pointer": "fuga",
			"ports": []interface{}{
				map[string]interface{}{"protocol": "TCP", "port": 53},
				map[string]interface{}{"protocol": "UDP", "port": 53, "targetPort": 1053},
			},
		}

		err := k8sClient.Patch(context.TODO(), obj, client.Apply, &client.PatchOptions{
			FieldManager: "foo",
		})
		Expect(err).ShouldNot(HaveOccurred())

		foo := myService()
		err = k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "foo"}, foo)
		Expect(err).ShouldNot(HaveOccurred())

		s, ok, err := unstructured.NestedString(foo.Object, "spec", "string")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(s).To(Equal("hoge"))
		s, ok, err = unstructured.NestedString(foo.Object, "spec", "pointer")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(s).To(Equal("fuga"))

		ports, ok, err := unstructured.NestedSlice(foo.Object, "spec", "ports")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(ports).To(HaveLen(2))

		for _, i := range ports {
			port, ok := i.(map[string]interface{})
			Expect(ok).To(BeTrue())
			if cmp.Equal(port["protocol"], "TCP") {
				_, ok := port["targetPort"]
				Expect(ok).To(BeFalse())
			} else {
				Expect(cmp.Equal(port["protocol"], "UDP")).To(BeTrue())
				targetPort, ok, err := unstructured.NestedInt64(port, "targetPort")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ok).To(BeTrue())
				Expect(targetPort).To(BeEquivalentTo(1053))
			}
		}
	})

	It("can add a slice item by another fieldManager", func() {
		obj := myService()
		obj.SetName("foo")
		obj.SetNamespace("default")
		obj.UnstructuredContent()["spec"] = map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{"protocol": "TCP", "port": 443},
			},
		}

		err := k8sClient.Patch(context.TODO(), obj, client.Apply, &client.PatchOptions{
			FieldManager: "another",
		})
		Expect(err).ShouldNot(HaveOccurred())

		foo := myService()
		err = k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "foo"}, foo)
		Expect(err).ShouldNot(HaveOccurred())

		ports, ok, err := unstructured.NestedSlice(foo.Object, "spec", "ports")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(ports).To(HaveLen(3))
	})

	It("can remove slice items", func() {
		obj := myService()
		obj.SetName("foo")
		obj.SetNamespace("default")
		obj.UnstructuredContent()["spec"] = map[string]interface{}{}

		err := k8sClient.Patch(context.TODO(), obj, client.Apply, &client.PatchOptions{
			FieldManager: "foo",
		})
		Expect(err).ShouldNot(HaveOccurred())

		foo := myService()
		err = k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "foo"}, foo)
		Expect(err).ShouldNot(HaveOccurred())

		// fields cannot be removed by SSA.  Values are left, fields become unmanaged.
		_, ok, err := unstructured.NestedString(foo.Object, "spec", "string")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		_, ok, err = unstructured.NestedString(foo.Object, "spec", "pointer")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())

		ports, ok, err := unstructured.NestedSlice(foo.Object, "spec", "ports")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(ports).To(HaveLen(1))
	})
})
