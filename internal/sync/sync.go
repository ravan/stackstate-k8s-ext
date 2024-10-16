package sync

import (
	"fmt"
	"github.com/ravan/stackstate-client/stackstate/receiver"
	"github.com/ravan/stackstate-k8s-ext/internal/config"
	"github.com/ravan/stackstate-k8s-ext/internal/k8s"
	storagev1 "k8s.io/api/storage/v1"
	"log/slog"
	"sigs.k8s.io/yaml"
	"strings"
)

const (
	Source = "k8s-ext"
)

func Sync(conf *config.Kubernetes) (*receiver.Factory, error) {
	factory := receiver.NewFactory(Source, Source, conf.Cluster)
	client, err := k8s.NewClient(conf)
	if err != nil {
		return nil, err
	}
	storageClasses, err := client.GetStorageClasses()
	if err != nil {
		return nil, err
	}

	scLookup := make(map[string]*receiver.Component)

	for _, sc := range storageClasses.Items {
		c := mapStorageClass(&sc, factory)
		scLookup[c.Type.Name] = c
	}

	return factory, nil
}

func mapStorageClass(sc *storagev1.StorageClass, f *receiver.Factory) *receiver.Component {
	id := UrnStorageClass(sc.Name, f.Cluster)
	var c *receiver.Component
	if f.ComponentExists(id) {
		c = f.MustGetComponent(id)
	} else {
		c = f.MustNewComponent(id, sc.Name, "storageclass")
		c.Data.Layer = "Storage"
		c.Data.Domain = Source
		data := make(receiver.PropertyMap)
		c.AddLabelKey("cluster-name", f.Cluster)
		c.AddCustomPropertyMap(Source, &data)
		c.AddProperty("clusterNameIdentifier", fmt.Sprintf("urn:cluster:/kubernetes:%s", f.Cluster))
		c.SourceProperties = convertToSourceProperties(sc)
		c.SourceProperties["apiVersion"] = "storage.k8s.io/v1"
		c.SourceProperties["kind"] = "StorageClass"
	}
	return c
}

func convertToSourceProperties(obj any) map[string]any {
	bytes, err := yaml.Marshal(obj)
	result := make(map[string]interface{})
	if err != nil {
		slog.Error("failed to marshal object to yaml", slog.Any("error", obj))
		return result
	}
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		slog.Error("failed to unmarshal object to yaml", slog.Any("error", obj))
	}
	delete(result["metadata"].(map[string]any), "managedFields")
	return result

}

func UrnStorageClass(name, cluster string) string {
	return fmt.Sprintf("urn:kubernetes:%s:storageclass:/%s", cluster, strings.ToLower(name))
}

func UrnPVC(name, namespace, cluster string) string {
	return fmt.Sprintf("urn:kubernetes:%s:%s:persistent-volume-claim/%s", cluster, namespace, strings.ToLower(name))
}
