package k8s

import (
	"context"
	"fmt"

	"embed"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

/*
	TODO: Service Protocol needs to be dynamic
		  Namespace needs to be dynamic - currently hard coded.
*/

var (
	//go:embed helx-apps/*
	helxFs embed.FS
)

type Deployment struct {
	DName      string            // AppName + UserSession Slice
	DNamespace string            // CurrentNamespace
	DLabels    map[string]string // user-app=DName
	AName      string            // AppName
	AImage     string            // AppImage
	APort      int               // AppPort
}

func int32Ptr(i int32) *int32 { return &i }

// This function will govern the Deployments.
// And a struct containing pertinent information will be populated and
// passed to CreateDeployment() for creation.
func (d *Deployment) CreateDeployment() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1().Deployments("appstore-system")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.DName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: d.DLabels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.DLabels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  d.AName,
							Image: d.AImage,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: int32(d.APort),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func (d *Deployment) CreateService() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	servicesClient := clientset.CoreV1().Services("appstore-system")
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.DName + "-service",
			Namespace: "appstore-system",
			Labels: map[string]string{
				"app": d.DName,
			},
		},
		Spec: apiv1.ServiceSpec{
			// TODO: Make this prettier and Protocol not dynamic
			Ports: []apiv1.ServicePort{
				{
					Name:     d.AName,
					Port:     int32(d.APort),
					Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector:  map[string]string{},
			ClusterIP: "",
		},
	}

	// Create Service
	fmt.Println("Creating service...")
	result, err := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func ListDeployment() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments("appstore-system")
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func DeleteDeployment(appname string) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments("appstore-system")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), appname, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

// func CreateDeploymentFromFile(appId string) {
// 	// Deploy from a manifest file.
// 	// https://medium.com/@harshjniitr/reading-and-writing-k8s-resource-as-yaml-in-golang-81dc8c7ea800
// 	decode := scheme.Codecs.UniversalDeserializer().Decode
// 	appFile := fmt.Sprintf("helx-apps/%s.yml", appId)
// 	stream, err := helxFs.ReadFile(appFile)
// 	if err != nil {
// 		fmt.Println(err, appFile)
// 	}
// 	obj, gKV, _ := decode(stream, nil, nil)
// 	deployment := &appsv1.Deployment{}
// 	if gKV.Kind == "Deployment" {
// 		deployment = obj.(*appsv1.Deployment)
// 	}

// 	// creates the in-cluster config
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// creates the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// Specify namespace in Deployments object "appstore-system"
// 	deploymentsClient := clientset.AppsV1().Deployments("appstore-system")

// 	// Create Deployment
// 	fmt.Println("Creating deployment...")
// 	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
// }
