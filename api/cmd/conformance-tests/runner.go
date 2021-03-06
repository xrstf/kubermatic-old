package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/onsi/ginkgo/reporters"
	"go.uber.org/zap"

	kubermaticapiv1 "github.com/kubermatic/kubermatic/api/pkg/api/v1"
	clusterclient "github.com/kubermatic/kubermatic/api/pkg/cluster/client"
	kubermaticv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	kubermaticv1helper "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1/helper"
	"github.com/kubermatic/kubermatic/api/pkg/provider"
	"github.com/kubermatic/kubermatic/api/pkg/resources"
	apiclient "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/client"
	projectclient "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/client/project"
	apimodels "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/models"
	clusterv1alpha1 "github.com/kubermatic/machine-controller/pkg/apis/cluster/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	utilerror "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func podIsReady(p *corev1.Pod) bool {
	for _, c := range p.Status.Conditions {
		if c.Type == corev1.PodReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

type testScenario interface {
	Name() string
	Cluster(secrets secrets) *apimodels.CreateClusterSpec
	NodeDeployments(num int, secrets secrets) ([]apimodels.NodeDeployment, error)
	OS() apimodels.OperatingSystemSpec
}

func newRunner(scenarios []testScenario, opts *Opts, log *zap.SugaredLogger) *testRunner {
	return &testRunner{
		scenarios:                    scenarios,
		controlPlaneReadyWaitTimeout: opts.controlPlaneReadyWaitTimeout,
		deleteClusterAfterTests:      opts.deleteClusterAfterTests,
		secrets:                      opts.secrets,
		namePrefix:                   opts.namePrefix,
		clusterClientProvider:        opts.clusterClientProvider,
		seed:                         opts.seed,
		seedRestConfig:               opts.seedRestConfig,
		nodeCount:                    opts.nodeCount,
		repoRoot:                     opts.repoRoot,
		reportsRoot:                  opts.reportsRoot,
		clusterParallelCount:         opts.clusterParallelCount,
		PublicKeys:                   opts.publicKeys,
		workerName:                   opts.workerName,
		homeDir:                      opts.homeDir,
		seedClusterClient:            opts.seedClusterClient,
		seedGeneratedClient:          opts.seedGeneratedClient,
		log:                          log,
		existingClusterLabel:         opts.existingClusterLabel,
		openshift:                    opts.openshift,
		openshiftPullSecret:          opts.openshiftPullSecret,
		printGinkoLogs:               opts.printGinkoLogs,
		onlyTestCreation:             opts.onlyTestCreation,
		pspEnabled:                   opts.pspEnabled,
		kubermatcProjectID:           opts.kubermatcProjectID,
		kubermaticClient:             opts.kubermaticClient,
		kubermaticAuthenticator:      opts.kubermaticAuthenticator,
	}
}

type testRunner struct {
	ctx                 context.Context
	scenarios           []testScenario
	secrets             secrets
	namePrefix          string
	repoRoot            string
	reportsRoot         string
	PublicKeys          [][]byte
	workerName          string
	homeDir             string
	log                 *zap.SugaredLogger
	openshift           bool
	openshiftPullSecret string
	printGinkoLogs      bool
	onlyTestCreation    bool
	pspEnabled          bool

	controlPlaneReadyWaitTimeout time.Duration
	deleteClusterAfterTests      bool
	nodeCount                    int
	clusterParallelCount         int

	seedClusterClient     ctrlruntimeclient.Client
	seedGeneratedClient   kubernetes.Interface
	clusterClientProvider *clusterclient.Provider
	seed                  *kubermaticv1.Seed
	seedRestConfig        *rest.Config

	// The label to use to select an existing cluster to test against instead of
	// creating a new one
	existingClusterLabel string

	kubermatcProjectID      string
	kubermaticClient        *apiclient.Kubermatic
	kubermaticAuthenticator runtime.ClientAuthInfoWriter
}

type testResult struct {
	report   *reporters.JUnitTestSuite
	err      error
	scenario testScenario
}

func (t *testResult) Passed() bool {
	if t.err != nil {
		return false
	}

	if t.report == nil {
		return false
	}

	if len(t.report.TestCases) == 0 {
		return false
	}

	if t.report.Errors > 0 || t.report.Failures > 0 {
		return false
	}

	return true
}

func (r *testRunner) worker(id int, scenarios <-chan testScenario, results chan<- testResult) {
	for s := range scenarios {
		scenarioLog := r.log.With("scenario", s.Name(), "worker", id)
		scenarioLog.Info("Starting to test scenario...")

		report, err := r.executeScenario(scenarioLog, s)
		res := testResult{
			report:   report,
			scenario: s,
			err:      err,
		}
		if err != nil {
			scenarioLog.Infof("Finished with error: %v", err)
		} else {
			scenarioLog.Info("Finished")
		}

		results <- res
	}
}

func (r *testRunner) Run() error {
	scenariosCh := make(chan testScenario, len(r.scenarios))
	resultsCh := make(chan testResult, len(r.scenarios))

	r.log.Info("Test suite:")
	for _, scenario := range r.scenarios {
		r.log.Info(scenario.Name())
		scenariosCh <- scenario
	}
	r.log.Info(fmt.Sprintf("Total: %d tests", len(r.scenarios)))

	for i := 1; i <= r.clusterParallelCount; i++ {
		go r.worker(i, scenariosCh, resultsCh)
	}

	close(scenariosCh)

	var results []testResult
	for range r.scenarios {
		results = append(results, <-resultsCh)
		r.log.Infof("Finished %d/%d test cases", len(results), len(r.scenarios))
	}

	overallResultBuf := &bytes.Buffer{}
	hadFailure := false
	for _, result := range results {
		prefix := "PASS"
		if !result.Passed() {
			prefix = "FAIL"
			hadFailure = true
		}
		scenarioResultMsg := fmt.Sprintf("[%s] - %s", prefix, result.scenario.Name())
		if result.err != nil {
			scenarioResultMsg = fmt.Sprintf("%s : %v", scenarioResultMsg, result.err)
		}

		fmt.Fprintln(overallResultBuf, scenarioResultMsg)
		if result.report != nil {
			printDetailedReport(result.report)
		}
	}

	fmt.Println("========================== RESULT ===========================")
	fmt.Println(overallResultBuf.String())

	if hadFailure {
		return errors.New("some tests failed")
	}
	return nil
}

func (r *testRunner) executeScenario(log *zap.SugaredLogger, scenario testScenario) (*reporters.JUnitTestSuite, error) {
	var err error
	var cluster *kubermaticv1.Cluster

	report := &reporters.JUnitTestSuite{
		Name: scenario.Name(),
	}
	totalStart := time.Now()

	// We'll store the report there and all kinds of logs
	scenarioFolder := path.Join(r.reportsRoot, scenario.Name())
	if err := os.MkdirAll(scenarioFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create the scenario folder '%s': %v", scenarioFolder, err)
	}

	// We need the closure to defer the evaluation of the time.Since(totalStart) call
	defer func() { log.Infof("Finished testing cluster after %s", time.Since(totalStart)) }()
	// Always write junit to disk
	defer func() {
		report.Time = time.Since(totalStart).Seconds()
		b, err := xml.Marshal(report)
		if err != nil {
			log.Errorw("failed to marshal junit", zap.Error(err))
			return
		}
		if err := ioutil.WriteFile(path.Join(r.reportsRoot, fmt.Sprintf("junit.%s.xml", scenario.Name())), b, 0644); err != nil {
			log.Errorw("failed to write junit", zap.Error(err))
		}
	}()

	ctx := context.Background()
	if r.existingClusterLabel == "" && os.Getenv("KUBERMATIC_USE_EXISTING_CLUSTER") == "" {
		if err := junitReporterWrapper(
			"[Kubermatic] Create cluster",
			report,
			func() error {
				cluster, err = r.createCluster(log, scenario)
				return err
			}); err != nil {
			return report, fmt.Errorf("failed to create cluster: %v", err)
		}
	} else {
		log.Infow("Using existing cluster")
		selector, err := labels.Parse(r.existingClusterLabel)
		if err != nil {
			return nil, fmt.Errorf("failed to parse labelselector %q: %v", r.existingClusterLabel, err)
		}
		clusterList := &kubermaticv1.ClusterList{}
		listOptions := &ctrlruntimeclient.ListOptions{LabelSelector: selector}
		if err := r.seedClusterClient.List(ctx, clusterList, listOptions); err != nil {
			return nil, fmt.Errorf("failed to list clusters: %v", err)
		}
		if foundClusterNum := len(clusterList.Items); foundClusterNum != 1 {
			return nil, fmt.Errorf("expected to find exactly one existing cluster, but got %d", foundClusterNum)
		}
		cluster = &clusterList.Items[0]
	}
	clusterName := cluster.Name
	log = log.With("cluster", cluster.Name)

	if err := junitReporterWrapper(
		"[Kubermatic] Wait for successful reconciliation",
		report,
		func() error {
			return wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
				if err := r.seedClusterClient.Get(ctx, types.NamespacedName{Name: clusterName}, cluster); err != nil {
					log.Errorw("Failed to get cluster when waiting for successful reconciliation", zap.Error(err))
					return false, nil
				}

				missingConditions, success := kubermaticv1helper.ClusterReconciliationSuccessful(cluster)
				if len(missingConditions) > 0 {
					log.Infof("Waiting for the following conditions: %v", missingConditions)
				}
				return success, nil
			})
		},
	); err != nil {
		return report, fmt.Errorf("failed to wait for successful reconciliation: %v", err)
	}

	if err := r.executeTests(log, cluster, report, scenario); err != nil {
		return report, err
	}

	if !r.deleteClusterAfterTests {
		return report, nil
	}

	return report, r.deleteCluster(report, cluster, log)
}

func (r *testRunner) executeTests(
	log *zap.SugaredLogger,
	cluster *kubermaticv1.Cluster,
	report *reporters.JUnitTestSuite,
	scenario testScenario,
) error {

	// We must store the name here because the cluster object may be nil on error
	clusterName := cluster.Name

	// Print all controlplane logs to both make debugging easier and show issues
	// that didn't result in test failures.
	defer r.printAllControlPlaneLogs(log, clusterName)

	var err error

	if err := junitReporterWrapper(
		"[Kubermatic] Wait for controlplane",
		report,
		func() error {
			cluster, err = r.waitForControlPlane(log, clusterName)
			return err
		},
	); err != nil {
		return fmt.Errorf("failed waiting for control plane to become ready: %v", err)
	}

	if err := junitReporterWrapper(
		"[Kubermatic] Add LB and PV Finalizers",
		report,
		func() error {
			return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
				if err := r.seedClusterClient.Get(context.Background(), types.NamespacedName{Name: clusterName}, cluster); err != nil {
					return err
				}
				cluster.Finalizers = append(cluster.Finalizers,
					kubermaticapiv1.InClusterPVCleanupFinalizer,
					kubermaticapiv1.InClusterLBCleanupFinalizer,
				)
				return r.seedClusterClient.Update(context.Background(), cluster)
			})
		},
	); err != nil {
		return fmt.Errorf("failed to add PV and LB cleanup finalizers: %v", err)
	}

	providerName, err := provider.ClusterCloudProviderName(cluster.Spec.Cloud)
	if err != nil {
		return fmt.Errorf("failed to get cloud provider name from cluster: %v", err)
	}

	log = log.With(
		"cloud-provider", providerName,
		"version", cluster.Spec.Version,
	)

	_, exists := r.seed.Spec.Datacenters[cluster.Spec.Cloud.DatacenterName]
	if !exists {
		return fmt.Errorf("datacenter %q doesn't exist", cluster.Spec.Cloud.DatacenterName)
	}

	kubeconfigFilename, err := r.getKubeconfig(log, cluster)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %v", err)
	}
	log = log.With("kubeconfig", kubeconfigFilename)

	cloudConfigFilename, err := r.getCloudConfig(log, cluster)
	if err != nil {
		return fmt.Errorf("failed to get cloud config: %v", err)
	}

	userClusterClient, err := r.clusterClientProvider.GetClient(cluster)
	if err != nil {
		return fmt.Errorf("failed to get the client for the cluster: %v", err)
	}

	if err := junitReporterWrapper(
		"[Kubermatic] Create NodeDeployments",
		report,
		func() error {
			return r.createNodeDeployments(log, scenario, clusterName)
		},
	); err != nil {
		return fmt.Errorf("failed to setup nodes: %v", err)
	}

	defer logEventsForAllMachines(context.Background(), log, userClusterClient)
	defer logUserClusterPodEventsAndLogs(
		log,
		r.clusterClientProvider,
		cluster.DeepCopy(),
	)

	var overallTimeout = 10 * time.Minute
	var timeoutLeft time.Duration
	if cluster.IsOpenshift() {
		// Openshift installs a lot more during node provisioning, hence this may take longer
		overallTimeout += 5 * time.Minute
	}
	// The initialization of the external CCM is super slow
	if cluster.Spec.Features[kubermaticv1.ClusterFeatureExternalCloudProvider] {
		overallTimeout += 5 * time.Minute
	}
	// Packet is slower at provisioning the instances, presumably because those are actual
	// physical hosts.
	if cluster.Spec.Cloud.Packet != nil {
		overallTimeout += 5 * time.Minute
	}

	if err := junitReporterWrapper(
		"[Kubermatic] Wait for machines to get a node",
		report,
		func() error {
			var err error
			timeoutLeft, err = waitForMachinesToJoinCluster(log, userClusterClient, overallTimeout)
			return err
		},
	); err != nil {
		return fmt.Errorf("failed to wait for machines to get a node: %v", err)
	}

	if err := junitReporterWrapper(
		"[Kubermatic] Wait for nodes to be ready",
		report,
		func() error {
			// Getting ready just implies starting the CNI deamonset, so that should
			// be quick.
			var err error
			timeoutLeft, err = waitForNodesToBeReady(log, userClusterClient, timeoutLeft)
			return err
		},
	); err != nil {
		return fmt.Errorf("failed to wait for all nodes to be ready: %v", err)
	}

	if err := junitReporterWrapper(
		"[Kubermatic] Wait for Pods inside usercluster to be ready",
		report,
		func() error {
			return r.waitUntilAllPodsAreReady(log, userClusterClient, timeoutLeft)
		},
	); err != nil {
		return fmt.Errorf("failed to wait for all pods to get ready: %v", err)
	}

	if r.onlyTestCreation {
		return nil
	}

	if err := r.testCluster(
		log,
		scenario,
		cluster,
		userClusterClient,
		kubeconfigFilename,
		cloudConfigFilename,
		report,
	); err != nil {
		return fmt.Errorf("failed to test cluster: %v", err)
	}

	return nil
}

func (r *testRunner) deleteCluster(report *reporters.JUnitTestSuite, cluster *kubermaticv1.Cluster, log *zap.SugaredLogger) error {

	deleteParms := &projectclient.DeleteClusterParams{
		ProjectID: r.kubermatcProjectID,
		DC:        r.seed.Name,
	}
	deleteTimeout := 15 * time.Minute
	if cluster.Spec.Cloud.Azure != nil {
		// 15 Minutes are not enough for Azure
		deleteTimeout = 30 * time.Minute
	}
	deleteParms.SetTimeout(15 * time.Second)
	if err := junitReporterWrapper(
		"[Kubermatic] Delete cluster",
		report,
		func() error {
			var selector labels.Selector
			var err error
			if r.workerName != "" {
				selector, err = labels.Parse(fmt.Sprintf("worker-name=%s", r.workerName))
				if err != nil {
					return fmt.Errorf("failed to parse selector: %v", err)
				}
			}
			return wait.PollImmediate(5*time.Second, deleteTimeout, func() (bool, error) {
				clusterList := &kubermaticv1.ClusterList{}
				listOpts := &ctrlruntimeclient.ListOptions{LabelSelector: selector}
				if err := r.seedClusterClient.List(r.ctx, clusterList, listOpts); err != nil {
					log.Errorw("Listing clusters failed", zap.Error(err))
					return false, nil
				}
				// Success!
				if len(clusterList.Items) == 0 {
					return true, nil
				}
				// Should never happen
				if len(clusterList.Items) > 1 {
					return false, fmt.Errorf("expected to find zero or one cluster, got %d", len(clusterList.Items))
				}
				// Cluster is currently being deleted
				if clusterList.Items[0].DeletionTimestamp != nil {
					return false, nil
				}
				// Issue Delete call
				log.With("cluster", clusterList.Items[0].Name).Info("Issuing DELETE call for cluster")
				deleteParms.ClusterID = clusterList.Items[0].Name
				_, err := r.kubermaticClient.Project.DeleteCluster(deleteParms, r.kubermaticAuthenticator)
				log.Infow("Issued cluster delete call", zap.Error(errors.New(fmtSwaggerError(err))))
				return false, nil
			})
		},
	); err != nil {
		log.Errorw("failed to delete cluster", zap.Error(err))
		return err
	}

	return nil
}

func retryNAttempts(maxAttempts int, f func(attempt int) error) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = f(attempt)
		if err != nil {
			continue
		}
		return nil
	}
	return fmt.Errorf("function did not succeeded after %d attempts: %v", maxAttempts, err)
}

func (r *testRunner) testCluster(
	log *zap.SugaredLogger,
	scenario testScenario,
	cluster *kubermaticv1.Cluster,
	userClusterClient ctrlruntimeclient.Client,
	kubeconfigFilename string,
	cloudConfigFilename string,
	report *reporters.JUnitTestSuite,
) error {
	const maxTestAttempts = 3
	var err error
	log.Info("Starting to test cluster...")

	if r.openshift {
		// Openshift supports neither the conformance tests nor PVs/LBs yet :/
		return nil
	}

	ginkgoRuns, err := r.getGinkgoRuns(log, scenario, kubeconfigFilename, cloudConfigFilename, cluster)
	if err != nil {
		return fmt.Errorf("failed to get Ginkgo runs: %v", err)
	}
	for _, run := range ginkgoRuns {
		if err := junitReporterWrapper(
			fmt.Sprintf("[Ginkgo] Run ginkgo tests %q", run.name),
			report,
			func() error {
				ginkgoRes, err := r.executeGinkgoRunWithRetries(log, run, userClusterClient)
				if ginkgoRes != nil {
					// We append the report from Ginkgo to our scenario wide report
					appendReport(report, ginkgoRes.report)
				}
				return err
			},
		); err != nil {
			// We still wan't to run potential next runs
			continue
		}

	}

	// Do a simple PVC test - with retries
	if supportsStorage(cluster) {
		if err := junitReporterWrapper(
			"[Kubermatic] [CloudProvider] Test PersistentVolumes",
			report,
			func() error {
				return retryNAttempts(maxTestAttempts,
					func(attempt int) error { return r.testPVC(log, userClusterClient, attempt) })
			},
		); err != nil {
			log.Errorf("Failed to verify that PVC's work: %v", err)
		}
	}

	// Do a simple LB test - with retries
	if supportsLBs(cluster) {
		if err := junitReporterWrapper(
			"[Kubermatic] [CloudProvider] Test LoadBalancers",
			report,
			func() error {
				return retryNAttempts(maxTestAttempts,
					func(attempt int) error { return r.testLB(log, userClusterClient, attempt) })
			},
		); err != nil {
			log.Errorf("Failed to verify that LB's work: %v", err)
		}
	}

	// Do user cluster RBAC controller test - with retries
	if err := junitReporterWrapper(
		"[Kubermatic] Test user cluster RBAC controller",
		report,
		func() error {
			return retryNAttempts(maxTestAttempts,
				func(attempt int) error {
					return r.testUserclusterControllerRBAC(log, cluster, userClusterClient, r.seedClusterClient)
				})
		}); err != nil {
		log.Errorf("Failed to verify that user cluster RBAC controller work: %v", err)
	}

	// Do prometheus metrics available test - with retries
	if err := junitReporterWrapper(
		"[Kubermatic] Test prometheus metrics availability", report, func() error {
			return retryNAttempts(maxTestAttempts, func(attempt int) error {
				return r.testUserClusterMetrics(log, cluster, r.seedClusterClient)
			})
		}); err != nil {
		log.Errorf("Failed to verify that prometheus metrics are available: %v", err)
	}

	return nil
}

// executeGinkgoRunWithRetries executes the passed GinkgoRun and retries if it failed hard(Failed to execute the Ginkgo binary for example)
// Or if the JUnit report from Ginkgo contains failed tests.
// Only if Ginkgo failed hard, an error will be returned. If some tests still failed after retrying the run, the report will reflect that.
func (r *testRunner) executeGinkgoRunWithRetries(log *zap.SugaredLogger, run *ginkgoRun, client ctrlruntimeclient.Client) (ginkgoRes *ginkgoResult, err error) {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ginkgoRes, err = executeGinkgoRun(log, run, client)
		if err != nil {
			// Something critical happened and we don't have a valid result
			log.Errorf("failed to execute the Ginkgo run '%s': %v", run.name, err)
			continue
		}

		if ginkgoRes.report.Errors > 0 || ginkgoRes.report.Failures > 0 {
			msg := fmt.Sprintf("Ginkgo run '%s' had failed tests.", run.name)
			if attempt < maxAttempts {
				msg = fmt.Sprintf("%s. Retrying...", msg)
			}
			log.Info(msg)
			if r.printGinkoLogs {
				if err := printFileUnbuffered(ginkgoRes.logfile); err != nil {
					log.Infof("error printing ginkgo logfile: %v", err)
				}
				log.Info("Successfully printed logfile")
			}
			continue
		}

		// Ginkgo run successfully and no test failed
		return ginkgoRes, err
	}

	return ginkgoRes, err
}

func (r *testRunner) createNodeDeployments(log *zap.SugaredLogger, scenario testScenario, clusterName string) error {

	var existingReplicas int
	log.Info("Getting existing NodeDeployments")
	nodeDeploymentGetParams := &projectclient.ListNodeDeploymentsParams{
		ProjectID: r.kubermatcProjectID,
		ClusterID: clusterName,
		DC:        r.seed.Name,
	}
	nodeDeploymentGetParams.SetTimeout(15 * time.Second)
	if err := wait.PollImmediate(10*time.Second, time.Minute, func() (bool, error) {
		resp, err := r.kubermaticClient.Project.ListNodeDeployments(nodeDeploymentGetParams, r.kubermaticAuthenticator)
		if err != nil {
			log.Errorw("Failed to get existing NodeDeployments", zap.Error(errors.New(fmtSwaggerError(err))))
			return false, nil
		}
		for _, nodeDeployment := range resp.Payload {
			existingReplicas += int(*nodeDeployment.Spec.Replicas)
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("failed to get existing NodeDeployments: %v", err)
	}
	log.Infof("Found %d pre-existing node replicas", existingReplicas)

	nodeCount := r.nodeCount - existingReplicas
	if nodeCount < 0 {
		return fmt.Errorf("found %d existing replicas and want %d, scaledown not supported", existingReplicas, r.nodeCount)
	}
	if nodeCount == 0 {
		return nil
	}

	log.Info("Creating NodeDeployments via kubermatic API")
	var nodeDeployments []apimodels.NodeDeployment
	var err error
	if err := wait.PollImmediate(10*time.Second, time.Minute, func() (bool, error) {
		nodeDeployments, err = scenario.NodeDeployments(nodeCount, r.secrets)
		if err != nil {
			log.Info("Getting NodeDeployments from scenario failed", zap.Error(err))
			return false, nil
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("didn't get NodeDeployments from scenario within a minute: %v", err)
	}

	for _, nd := range nodeDeployments {
		params := &projectclient.CreateNodeDeploymentParams{
			ProjectID: r.kubermatcProjectID,
			ClusterID: clusterName,
			DC:        r.seed.Name,
			Body:      &nd,
		}
		params.SetTimeout(15 * time.Second)

		if err := retryNAttempts(defaultAPIRetries, func(attempt int) error {
			if _, err := r.kubermaticClient.Project.CreateNodeDeployment(params, r.kubermaticAuthenticator); err != nil {
				log.Warnf("[Attempt %d/%d] Failed to create NodeDeployment %s: %v. Retrying", attempt, defaultAPIRetries, nd.Name, fmtSwaggerError(err))
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to create NodeDeployment %s via kubermatic api after %d attempts: %q", nd.Name, defaultAPIRetries, fmtSwaggerError(err))
		}
	}

	log.Infof("Successfully created %d NodeDeployments via Kubermatic API", nodeCount)
	return nil
}

func (r *testRunner) getKubeconfig(log *zap.SugaredLogger, cluster *kubermaticv1.Cluster) (string, error) {
	log.Debug("Getting kubeconfig...")
	var kubeconfig []byte
	// Needed for Openshift where we have to create a SA and bindings inside the cluster
	// which can only be done after the APIServer is up and ready
	if err := wait.PollImmediate(time.Second, time.Minute, func() (bool, error) {
		var err error
		kubeconfig, err = r.clusterClientProvider.GetAdminKubeconfig(cluster)
		if err != nil {
			log.Debugw("Failed to get Kubeconfig", zap.Error(err))
			return false, nil
		}
		return true, nil
	}); err != nil {
		return "", fmt.Errorf("failed to wait for kubeconfig: %v", err)
	}
	filename := path.Join(r.homeDir, fmt.Sprintf("%s-kubeconfig", cluster.Name))
	if err := ioutil.WriteFile(filename, kubeconfig, 0644); err != nil {
		return "", fmt.Errorf("failed to write kubeconfig to %s: %v", filename, err)
	}

	log.Infof("Successfully wrote kubeconfig to %s", filename)
	return filename, nil
}

func (r *testRunner) getCloudConfig(log *zap.SugaredLogger, cluster *kubermaticv1.Cluster) (string, error) {
	log.Debug("Getting cloud-config...")

	var cmData string
	err := retryNAttempts(defaultAPIRetries, func(attempt int) error {
		cm := &corev1.ConfigMap{}
		name := types.NamespacedName{Namespace: cluster.Status.NamespaceName, Name: resources.CloudConfigConfigMapName}
		if err := r.seedClusterClient.Get(context.Background(), name, cm); err != nil {
			return fmt.Errorf("failed to load cloud-config: %v", err)
		}
		cmData = cm.Data["config"]
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to get cloud config ConfigMap: %v", err)
	}

	filename := path.Join(r.homeDir, fmt.Sprintf("%s-cloud-config", cluster.Name))
	if err := ioutil.WriteFile(filename, []byte(cmData), 0644); err != nil {
		return "", fmt.Errorf("failed to write cloud config: %v", err)
	}

	log.Infof("Successfully wrote cloud-config to %s", filename)
	return filename, nil
}

func (r *testRunner) createCluster(log *zap.SugaredLogger, scenario testScenario) (*kubermaticv1.Cluster, error) {
	log.Info("Creating cluster via kubermatic API")

	cluster := scenario.Cluster(r.secrets)
	if r.openshift {
		cluster.Cluster.Type = "openshift"
		cluster.Cluster.Spec.Openshift = &apimodels.Openshift{
			ImagePullSecret: r.openshiftPullSecret,
		}
	}
	// The cluster name must be unique per project.
	// We build up a understandable name with the various cli parameters & add a random string in the end to ensure
	// we really have a unique name
	if r.namePrefix != "" {
		cluster.Cluster.Name = r.namePrefix + "-"
	}
	if r.workerName != "" {
		cluster.Cluster.Name += r.workerName + "-"
	}
	cluster.Cluster.Name += scenario.Name() + "-"
	cluster.Cluster.Name += rand.String(8)

	cluster.Cluster.Spec.UsePodSecurityPolicyAdmissionPlugin = r.pspEnabled

	params := &projectclient.CreateClusterParams{
		ProjectID: r.kubermatcProjectID,
		DC:        r.seed.Name,
		Body:      cluster,
	}
	params.SetTimeout(15 * time.Second)

	crCluster := &kubermaticv1.Cluster{}
	var selector labels.Selector
	var err error
	if r.workerName != "" {
		selector, err = labels.Parse(fmt.Sprintf("worker-name=%s", r.workerName))
		if err != nil {
			return nil, fmt.Errorf("failed to parse selector: %v", err)
		}
	}

	var errs []error
	if err := wait.PollImmediate(5*time.Second, 45*time.Second, func() (bool, error) {
		// For some reason the cluster doesn't have the name we set via ID on creation
		clusterList := &kubermaticv1.ClusterList{}
		opts := &ctrlruntimeclient.ListOptions{LabelSelector: selector}
		if err := r.seedClusterClient.List(r.ctx, clusterList, opts); err != nil {
			return false, err
		}
		numFoundClusters := len(clusterList.Items)

		switch {
		case numFoundClusters < 1:
			if _, err := r.kubermaticClient.Project.CreateCluster(params, r.kubermaticAuthenticator); err != nil {
				// Log the error but don't return it, we want to retry
				err = errors.New(fmtSwaggerError(err))
				errs = append(errs, err)
				log.Errorf("failed to create cluster via kubermatic api: %q", err)
			} else {
				log.Info("Successfully created cluster via kubermatic api")
			}
			// Always return here, our clusterList is not up to date anymore
			return false, nil
		case numFoundClusters > 1:
			return false, fmt.Errorf("had more than one cluster (%d) with our worker-name, how is this possible?! ", numFoundClusters)
		default:
			crCluster = &clusterList.Items[0]
			return true, err
		}
	}); err != nil {
		errs = append(errs, err)
		return nil, fmt.Errorf("cluster creation failed: %v", utilerror.NewAggregate(errs))
	}

	log.Info("Successfully created cluster via Kubermatic API")
	return crCluster, nil
}

func (r *testRunner) waitForControlPlane(log *zap.SugaredLogger, clusterName string) (*kubermaticv1.Cluster, error) {
	log.Debug("Waiting for control plane to become ready...")
	started := time.Now()
	namespacedClusterName := types.NamespacedName{Name: clusterName}

	err := wait.Poll(controlPlaneReadyPollPeriod, r.controlPlaneReadyWaitTimeout, func() (done bool, err error) {
		newCluster := &kubermaticv1.Cluster{}

		if err := r.seedClusterClient.Get(context.Background(), namespacedClusterName, newCluster); err != nil {
			if kerrors.IsNotFound(err) {
				return false, nil
			}
		}
		// Check for this first, because otherwise we instantly return as the cluster-controller did not
		// create any pods yet
		if !newCluster.Status.ExtendedHealth.AllHealthy() {
			return false, nil
		}

		controlPlanePods := &corev1.PodList{}
		if err := r.seedClusterClient.List(
			context.Background(),
			controlPlanePods,
			&ctrlruntimeclient.ListOptions{Namespace: newCluster.Status.NamespaceName},
		); err != nil {
			return false, fmt.Errorf("failed to list controlplane pods: %v", err)
		}
		for _, pod := range controlPlanePods.Items {
			if !podIsReady(&pod) {
				return false, nil
			}
		}

		return true, nil
	})
	// Timeout or other error
	if err != nil {
		return nil, err
	}

	// Get copy of latest version
	cluster := &kubermaticv1.Cluster{}
	if err := r.seedClusterClient.Get(context.Background(), namespacedClusterName, cluster); err != nil {
		return nil, err
	}

	log.Debugf("Control plane became ready after %.2f seconds", time.Since(started).Seconds())
	return cluster, nil
}

func (r *testRunner) waitUntilAllPodsAreReady(log *zap.SugaredLogger, userClusterClient ctrlruntimeclient.Client, timeout time.Duration) error {
	log.Debug("Waiting for all pods to be ready...")
	started := time.Now()

	err := wait.Poll(defaultUserClusterPollInterval, timeout, func() (done bool, err error) {
		podList := &corev1.PodList{}
		if err := userClusterClient.List(context.Background(), podList); err != nil {
			log.Warnf("failed to load pod list while waiting until all pods are running: %v", err)
			return false, nil
		}

		for _, pod := range podList.Items {
			if !podIsReady(&pod) {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	log.Debugf("All pods became ready after %.2f seconds", time.Since(started).Seconds())
	return nil
}

type ginkgoResult struct {
	logfile  string
	report   *reporters.JUnitTestSuite
	duration time.Duration
}

const (
	argSeparator = ` \
    `
)

type ginkgoRun struct {
	name       string
	cmd        *exec.Cmd
	reportsDir string
	timeout    time.Duration
}

func (r *testRunner) getGinkgoRuns(
	log *zap.SugaredLogger,
	scenario testScenario,
	kubeconfigFilename,
	cloudConfigFilename string,
	cluster *kubermaticv1.Cluster,
) ([]*ginkgoRun, error) {
	kubeconfigFilename = path.Clean(kubeconfigFilename)
	repoRoot := path.Clean(r.repoRoot)
	MajorMinor := fmt.Sprintf("%d.%d", cluster.Spec.Version.Major(), cluster.Spec.Version.Minor())

	nodeNumberTotal := int32(r.nodeCount)

	ginkgoSkipParallel := `\[Serial\]`
	if minor := cluster.Spec.Version.Minor(); minor == 16 || minor == 17 {
		// These require the nodes NodePort to be available from the tester, which is not the case for us.
		// TODO: Maybe add an option to allow the NodePorts in the SecurityGroup?
		ginkgoSkipParallel += "|Services should be able to change the type from ExternalName to NodePort|Services should be able to create a functioning NodePort service"
	}

	runs := []struct {
		name          string
		ginkgoFocus   string
		ginkgoSkip    string
		parallelTests int
		timeout       time.Duration
	}{
		{
			name:          "parallel",
			ginkgoFocus:   `\[Conformance\]`,
			ginkgoSkip:    ginkgoSkipParallel,
			parallelTests: int(nodeNumberTotal) * 10,
			timeout:       15 * time.Minute,
		},
		{
			name:          "serial",
			ginkgoFocus:   `\[Serial\].*\[Conformance\]`,
			ginkgoSkip:    `should not cause race condition when used for configmap`,
			parallelTests: 1,
			timeout:       10 * time.Minute,
		},
	}
	versionRoot := path.Join(repoRoot, MajorMinor)
	binRoot := path.Join(versionRoot, "/platforms/linux/amd64")
	var ginkgoRuns []*ginkgoRun
	for _, run := range runs {

		reportsDir := path.Join("/tmp", scenario.Name(), run.name)
		env := []string{
			fmt.Sprintf("HOME=%s", r.homeDir),
			fmt.Sprintf("AWS_SSH_KEY=%s", path.Join(r.homeDir, ".ssh", "google_compute_engine")),
			fmt.Sprintf("LOCAL_SSH_KEY=%s", path.Join(r.homeDir, ".ssh", "google_compute_engine")),
			fmt.Sprintf("KUBE_SSH_KEY=%s", path.Join(r.homeDir, ".ssh", "google_compute_engine")),
		}

		args := []string{
			"-progress",
			fmt.Sprintf("-nodes=%d", run.parallelTests),
			"-noColor=true",
			"-flakeAttempts=2",
			fmt.Sprintf(`-focus=%s`, run.ginkgoFocus),
			fmt.Sprintf(`-skip=%s`, run.ginkgoSkip),
			path.Join(binRoot, "e2e.test"),
			"--",
			"--disable-log-dump",
			fmt.Sprintf("--repo-root=%s", versionRoot),
			fmt.Sprintf("--report-dir=%s", reportsDir),
			fmt.Sprintf("--report-prefix=%s", run.name),
			fmt.Sprintf("--kubectl-path=%s", path.Join(binRoot, "kubectl")),
			fmt.Sprintf("--kubeconfig=%s", kubeconfigFilename),
			fmt.Sprintf("--num-nodes=%d", nodeNumberTotal),
			fmt.Sprintf("--cloud-config-file=%s", cloudConfigFilename),
		}

		args = append(args, "--provider=local")

		osSpec := scenario.OS()
		switch {
		case osSpec.Ubuntu != nil:
			args = append(args, "--node-os-distro=ubuntu")
			env = append(env, "KUBE_SSH_USER=ubuntu")
		case osSpec.Centos != nil:
			args = append(args, "--node-os-distro=centos")
			env = append(env, "KUBE_SSH_USER=centos")
		case osSpec.ContainerLinux != nil:
			args = append(args, "--node-os-distro=coreos")
			env = append(env, "KUBE_SSH_USER=core")
		}

		cmd := exec.Command(path.Join(binRoot, "ginkgo"), args...)
		cmd.Env = env

		ginkgoRuns = append(ginkgoRuns, &ginkgoRun{
			name:       run.name,
			cmd:        cmd,
			reportsDir: reportsDir,
			timeout:    run.timeout,
		})
	}

	return ginkgoRuns, nil
}

func executeGinkgoRun(parentLog *zap.SugaredLogger, run *ginkgoRun, client ctrlruntimeclient.Client) (*ginkgoResult, error) {
	started := time.Now()
	log := parentLog.With("reports-dir", run.reportsDir)

	if err := deleteAllNonDefaultNamespaces(log, client); err != nil {
		return nil, fmt.Errorf("failed to cleanup namespaces before the Ginkgo run: %v", err)
	}

	// We're clearing up the temp dir on every run
	if err := os.RemoveAll(run.reportsDir); err != nil {
		log.Errorf("failed to remove temporary reports directory: %v", err)
	}
	if err := os.MkdirAll(run.reportsDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create temporary reports directory: %v", err)
	}

	// Make sure we write to a file instead of a byte buffer as the logs are pretty big
	file, err := ioutil.TempFile("/tmp", run.name+"-log")
	if err != nil {
		return nil, fmt.Errorf("failed to open logfile: %v", err)
	}
	defer file.Close()
	log = log.With("ginkgo-log", file.Name())

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	ctx, cancel := context.WithTimeout(context.Background(), run.timeout)
	defer cancel()

	// Copy the command as we cannot execute a command twice
	cmd := exec.CommandContext(ctx, "")
	cmd.Path = run.cmd.Path
	cmd.Args = run.cmd.Args
	cmd.Env = run.cmd.Env
	cmd.Dir = run.cmd.Dir
	cmd.ExtraFiles = run.cmd.ExtraFiles
	if _, err := writer.Write([]byte(strings.Join(cmd.Args, argSeparator))); err != nil {
		return nil, fmt.Errorf("failed to write command to log: %v", err)
	}

	log.Debugf("Starting Ginkgo run '%s'...", run.name)

	// Flush to disk so we can actually watch logs
	stopCh := make(chan struct{}, 1)
	defer close(stopCh)
	go wait.Until(func() {
		if err := writer.Flush(); err != nil {
			log.Warnf("failed to flush log writer: %v", err)
		}
		if err := file.Sync(); err != nil {
			log.Warnf("failed to sync log file: %v", err)
		}
	}, 1*time.Second, stopCh)

	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Debugf("Ginkgo exited with a non 0 return code: %v", exitErr)
		} else {
			return nil, fmt.Errorf("ginkgo failed to start: %T %v", err, err)
		}
	}

	// When running ginkgo in parallel, each ginkgo worker creates a own report, thus we must combine them
	combinedReport, err := collectReports(run.name, run.reportsDir)
	if err != nil {
		return nil, err
	}

	// If we have no junit files, we cannot return a valid report
	if len(combinedReport.TestCases) == 0 {
		return nil, errors.New("ginkgo report is empty. It seems no tests where executed")
	}

	combinedReport.Time = time.Since(started).Seconds()

	log.Debugf("Ginkgo run '%s' took %s", run.name, time.Since(started))
	return &ginkgoResult{
		logfile:  file.Name(),
		report:   combinedReport,
		duration: time.Since(started),
	}, nil
}

func supportsStorage(cluster *kubermaticv1.Cluster) bool {
	return cluster.Spec.Cloud.Openstack != nil ||
		cluster.Spec.Cloud.Azure != nil ||
		cluster.Spec.Cloud.AWS != nil ||
		cluster.Spec.Cloud.VSphere != nil ||
		cluster.Spec.Cloud.GCP != nil

	// Currently broken, see https://github.com/kubermatic/kubermatic/issues/3312
	//cluster.Spec.Cloud.Hetzner != nil
}

func supportsLBs(cluster *kubermaticv1.Cluster) bool {
	return cluster.Spec.Cloud.Azure != nil ||
		cluster.Spec.Cloud.AWS != nil ||
		cluster.Spec.Cloud.GCP != nil
}

func (r *testRunner) printAllControlPlaneLogs(log *zap.SugaredLogger, clusterName string) {
	log.Info("Printing controlplane logs")
	cluster := &kubermaticv1.Cluster{}
	ctx := context.Background()
	if err := r.seedClusterClient.Get(ctx, types.NamespacedName{Name: clusterName}, cluster); err != nil {
		log.Errorw("Failed to get cluster", zap.Error(err))
		return
	}

	clusterHealthStatus, err := json.Marshal(cluster.Status.ExtendedHealth)
	if err != nil {
		log.Errorw("Failed to marshal cluster health status", zap.Error(err))
	} else {
		fmt.Printf("ClusterHealthStatus: '%s'\n", clusterHealthStatus)
	}

	log.Infow("Logging events for cluster")
	if err := logEventsObject(ctx, log, r.seedClusterClient, "default", cluster.UID); err != nil {
		log.Errorw("Failed to log cluster events", zap.Error(err))
	}

	if err := printEventsAndLogsForAllPods(
		ctx,
		log,
		r.seedClusterClient,
		r.seedGeneratedClient,
		cluster.Status.NamespaceName,
	); err != nil {
		log.Errorw("Failed to print events and logs of pods", zap.Error(err))
	}
}

// waitForMachinesToJoinCluster waits for machines to join the cluster. It does so by checking
// if the machines have a nodeRef. It does not check if the nodeRef is valid.
// All errors are swallowed, only the timeout error is returned.
func waitForMachinesToJoinCluster(log *zap.SugaredLogger, client ctrlruntimeclient.Client, timeout time.Duration) (time.Duration, error) {
	startTime := time.Now()
	err := wait.Poll(10*time.Second, timeout, func() (bool, error) {
		machineList := &clusterv1alpha1.MachineList{}
		if err := client.List(context.Background(), machineList); err != nil {
			log.Warnw("Failed to list machines", zap.Error(err))
			return false, nil
		}
		for _, machine := range machineList.Items {
			if !machineHasNodeRef(machine) {
				log.Infow("Machine has no nodeRef yet", "machine", machine.Name)
				return false, nil
			}
		}
		log.Infow("All machines got a Node", "duration-in-seconds", time.Since(startTime).Seconds())
		return true, nil
	})
	return timeout - time.Since(startTime), err
}

func machineHasNodeRef(machine clusterv1alpha1.Machine) bool {
	return machine.Status.NodeRef != nil && machine.Status.NodeRef.Name != ""
}

// WaitForNodesToBeReady waits for all nodes to be ready. It does so by checking the Nodes "Ready"
// condition. It swallows all errors except for the timeout.
func waitForNodesToBeReady(log *zap.SugaredLogger, client ctrlruntimeclient.Client, timeout time.Duration) (time.Duration, error) {
	startTime := time.Now()
	err := wait.Poll(10*time.Second, timeout, func() (bool, error) {
		nodeList := &corev1.NodeList{}
		if err := client.List(context.Background(), nodeList); err != nil {
			log.Warnw("Failed to list nodes", zap.Error(err))
			return false, nil
		}
		for _, node := range nodeList.Items {
			if !nodeIsReady(node) {
				log.Infow("Node is not ready", "node", node.Name)
				return false, nil
			}
		}
		log.Infow("All nodes got ready", "duration-in-seconds", time.Since(startTime).Seconds())
		return true, nil
	})
	return timeout - time.Since(startTime), err
}

func nodeIsReady(node corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func printFileUnbuffered(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return printUnbuffered(fd)
}

// printUnbuffered uses io.Copy to print data to stdout.
// It should be used for all bigger logs, to avoid buffering
// them in memory and getting oom killed because of that.
func printUnbuffered(src io.Reader) error {
	_, err := io.Copy(os.Stdout, src)
	return err
}

// fmtSwaggerError works around the Error() implementration generated by swagger
// which prints only a pointer to the body but we want to see the actual content of the body.
// to fix this we can either type assert for each request type or naively use json
func fmtSwaggerError(err error) string {
	rawData, newErr := json.Marshal(err)
	if newErr != nil {
		return fmt.Sprintf("failed to marshal response(%v): %v", err, newErr)
	}
	return string(rawData)
}

// junitReporterWrapper is a convenience func to get junit results for a step
// It will create a report, append it to the passed in testsuite and propagate
// the error of the executor back up
// TODO: Should we add optional retrying here to limit the amount of wrappers we need?
func junitReporterWrapper(
	testCaseName string,
	report *reporters.JUnitTestSuite,
	executor func() error,
	extraErrOutputFn ...func() string,
) error {
	junitTestCase := reporters.JUnitTestCase{
		Name:      testCaseName,
		ClassName: testCaseName,
	}

	startTime := time.Now()
	err := executor()
	junitTestCase.Time = time.Since(startTime).Seconds()
	if err != nil {
		junitTestCase.FailureMessage = &reporters.JUnitFailureMessage{Message: err.Error()}
		report.Failures++
		for _, extraOut := range extraErrOutputFn {
			extraOutString := extraOut()
			err = fmt.Errorf("%v\n%s", err, extraOutString)
			junitTestCase.FailureMessage.Message += "\n" + extraOutString
		}
	}

	report.TestCases = append(report.TestCases, junitTestCase)
	report.Tests++

	return err
}

// printEvents and logs for all pods. Include ready pods, because they may still contain useful information.
func printEventsAndLogsForAllPods(
	ctx context.Context,
	log *zap.SugaredLogger,
	client ctrlruntimeclient.Client,
	k8sclient kubernetes.Interface,
	namespace string,
) error {
	log.Infow("Printing logs for all pods", "namespace", namespace)

	pods := &corev1.PodList{}
	if err := client.List(ctx, pods, ctrlruntimeclient.InNamespace(namespace)); err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	var errs []error
	for _, pod := range pods.Items {
		log := log.With("pod", pod.Name)
		if !podIsReady(&pod) {
			log.Error("Pod is not ready")
		}
		log.Info("Logging events for pod")
		if err := logEventsObject(ctx, log, client, pod.Namespace, pod.UID); err != nil {
			log.Errorw("Failed to log events for pod", zap.Error(err))
			errs = append(errs, err)
		}
		log.Info("Printing logs for pod")
		if err := printLogsForPod(log, k8sclient, &pod); err != nil {
			log.Errorw("Failed to print logs for pod", zap.Error(utilerror.NewAggregate(err)))
			errs = append(errs, err...)
		}
	}

	return utilerror.NewAggregate(errs)
}

func printLogsForPod(
	log *zap.SugaredLogger,
	k8sclient kubernetes.Interface,
	pod *corev1.Pod,
) []error {
	var errs []error
	for _, container := range pod.Spec.Containers {
		log.Infow("Printing logs for container", "container", container.Name)
		if err := printLogsForContainer(k8sclient, pod, container.Name); err != nil {
			log.Errorw(
				"Failed to print logs for container",
				"name", container.Name,
				zap.Error(err),
			)
			errs = append(errs, err)
		}
	}
	for _, initContainer := range pod.Spec.InitContainers {
		log.Infow("Printing logs for initContainer", "initContainer", initContainer.Name)
		if err := printLogsForContainer(k8sclient, pod, initContainer.Name); err != nil {
			log.Errorw(
				"Failed to print logs for container",
				"name", initContainer.Name,
				zap.Error(err),
			)
			errs = append(errs, err)
		}
	}
	return errs
}

func printLogsForContainer(client kubernetes.Interface, pod *corev1.Pod, containerName string) error {
	readCloser, err := client.
		CoreV1().
		Pods(pod.Namespace).
		GetLogs(pod.Name, &corev1.PodLogOptions{Container: containerName}).
		Stream()
	if err != nil {
		return err
	}
	defer readCloser.Close()
	return printUnbuffered(readCloser)
}

func logEventsForAllMachines(
	ctx context.Context,
	log *zap.SugaredLogger,
	client ctrlruntimeclient.Client,
) {
	machines := &clusterv1alpha1.MachineList{}
	if err := client.List(ctx, machines); err != nil {
		log.Errorw("Failed to list machines", zap.Error(err))
		return
	}

	for _, machine := range machines.Items {
		log.Infow("Logging events for machine", "name", machine.Name)
		if err := logEventsObject(ctx, log, client, machine.Namespace, machine.UID); err != nil {
			log.Errorw(
				"Failed to log events for machine",
				"name", machine.Name,
				"namespace", machine.Namespace,
				zap.Error(err),
			)
		}
	}
}

func logEventsObject(
	ctx context.Context,
	log *zap.SugaredLogger,
	client ctrlruntimeclient.Client,
	namespace string,
	uid types.UID,
) error {
	events := &corev1.EventList{}
	listOpts := &ctrlruntimeclient.ListOptions{
		Namespace:     namespace,
		FieldSelector: fields.OneTermEqualSelector("involvedObject.uid", string(uid)),
	}
	if err := client.List(ctx, events, listOpts); err != nil {
		return fmt.Errorf("failed to get events: %v", err)
	}

	for _, event := range events.Items {
		var msg string
		if event.Type == corev1.EventTypeWarning {
			// Make sure this gets highlighted
			msg = "ERROR"
		}
		log.Infow(
			msg,
			"EventType", event.Type,
			"Number", event.Count,
			"Reason", event.Reason,
			"Message", event.Message,
			"Source", event.Source.Component,
		)
	}
	return nil
}

func logUserClusterPodEventsAndLogs(
	log *zap.SugaredLogger,
	connProvider *clusterclient.Provider,
	cluster *kubermaticv1.Cluster,
) {
	log.Info("Attempting to log usercluster pod events and logs")
	cfg, err := connProvider.GetClientConfig(cluster)
	if err != nil {
		log.Errorw("Failed to get usercluster admin kubeconfig")
		return
	}
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Errorw("Failed to construct k8sClient for usercluster", zap.Error(err))
		return
	}
	client, err := connProvider.GetClient(cluster)
	if err != nil {
		log.Errorw("Failed to construct client for usercluster", zap.Error(err))
		return
	}
	if err := printEventsAndLogsForAllPods(
		context.Background(),
		log,
		client,
		k8sClient,
		"",
	); err != nil {
		log.Errorw("Failed to print events and logs for usercluster pods", zap.Error(err))
	}
}
