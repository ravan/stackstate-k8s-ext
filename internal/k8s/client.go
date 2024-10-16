package k8s

import (
	"context"
	"github.com/ravan/stackstate-k8s-ext/internal/config"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient(conf *config.Kubernetes) (*Client, error) {
	var c *rest.Config
	var err error
	if conf.InCluster {
		c, err = rest.InClusterConfig()
	} else {
		c, err = clientcmd.BuildConfigFromFlags("", conf.KubeConfig)
	}
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func (c *Client) GetStorageClasses() (*storagev1.StorageClassList, error) {
	scList, err := c.clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return scList, nil
}

func (c *Client) GetPVCs() (*corev1.PersistentVolumeClaimList, error) {
	pvcList, err := c.clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pvcList, nil
}
