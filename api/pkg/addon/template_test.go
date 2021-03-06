package addon

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"

	kubermaticv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	"github.com/kubermatic/kubermatic/api/pkg/resources"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// TestRenderAddons ensures that all our default addon manifests render
// properly given a variety of cluster configurations.
func TestRenderAddons(t *testing.T) {
	testRenderAddonsForOrchestrator(t, "kubernetes")
	testRenderAddonsForOrchestrator(t, "openshift")
}

func testRenderAddonsForOrchestrator(t *testing.T, orchestrator string) {
	clusterFiles, _ := filepath.Glob(fmt.Sprintf("testdata/cluster-%s-*", orchestrator))

	clusters := []kubermaticv1.Cluster{}
	for _, filename := range clusterFiles {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		cluster := kubermaticv1.Cluster{}
		err = yaml.Unmarshal(content, &cluster)
		if err != nil {
			t.Fatal(err)
		}

		clusters = append(clusters, cluster)
	}

	addonBasePath := "../../../addons"
	if orchestrator == "openshift" {
		addonBasePath = "../../../openshift_addons"
	}

	addonPaths, _ := filepath.Glob(filepath.Join(addonBasePath, "*"))

	addons := []kubermaticv1.Addon{}
	for _, addonPath := range addonPaths {
		addonPath, _ := filepath.Abs(addonPath)

		if stat, err := os.Stat(addonPath); err != nil || !stat.IsDir() {
			continue
		}

		addons = append(addons, kubermaticv1.Addon{
			ObjectMeta: metav1.ObjectMeta{
				Name: filepath.Base(addonPath),
			},
		})
	}

	log := zap.NewNop().Sugar()
	credentials := resources.Credentials{}
	variables := map[string]interface{}{
		"test": true,
	}

	for _, cluster := range clusters {
		for _, addon := range addons {
			data, err := NewTemplateData(&cluster, credentials, "kubeconfig", "1.2.3.4", "5.6.7.8", variables)
			if err != nil {
				t.Fatalf("Rendering %s addon %s for cluster %s failed: %v", orchestrator, addon.Name, cluster.Name, err)
			}

			path := filepath.Join(addonBasePath, addon.Name)

			_, err = ParseFromFolder(log, "", path, data)
			if err != nil {
				t.Fatalf("Rendering %s addon %s for cluster %s failed: %v", orchestrator, addon.Name, cluster.Name, err)
			}
		}
	}
}
