package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// GetActionConfig returns a configured action.Configuration for Helm operations
func (s *Server) GetActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	
	// Set up configFlags
	inCluster := os.Getenv("IN_CLUSTER") == "true"
	var configFlags *genericclioptions.ConfigFlags
	
	if inCluster {
		// Use in-cluster config
		configFlags = genericclioptions.NewConfigFlags(false)
	} else {
		// Use kubeconfig file
		configFlags = genericclioptions.NewConfigFlags(false)
		configFlags.KubeConfig = &s.KubeConfig
	}
	
	// Initialize with debug logging
	if err := actionConfig.Init(configFlags, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, fmt.Errorf("failed to initialize action configuration: %w", err)
	}

	return actionConfig, nil
}

// InstallChart installs a Helm chart with the given parameters
func (s *Server) InstallChart(params ChartInstallRequest) (*release.Release, error) {
	cfg, err := s.GetActionConfig(params.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize install action
	install := action.NewInstall(cfg)
	install.ReleaseName = params.ReleaseName
	install.Namespace = params.Namespace
	install.Wait = params.Wait
	install.CreateNamespace = params.CreateNamespace
	install.DryRun = params.DryRun
	install.ClientOnly = params.ClientOnly
	install.Description = params.Description

	// Set timeout if specified
	if params.Timeout > 0 {
		install.Timeout = time.Duration(params.Timeout) * time.Second
	}

	// Set chart path options if using repo
	if params.RepoURL != "" {
		install.ChartPathOptions.RepoURL = params.RepoURL
		if params.Version != "" {
			install.ChartPathOptions.Version = params.Version
		}
	}

	// Locate and load the chart
	chartPath, err := install.ChartPathOptions.LocateChart(params.ChartName, s.HelmEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	// Install the chart
	return install.Run(chart, params.Values)
}

// UpgradeChart upgrades an existing Helm chart release
func (s *Server) UpgradeChart(params ChartUpgradeRequest) (*release.Release, error) {
	cfg, err := s.GetActionConfig(params.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize upgrade action
	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = params.Namespace
	upgrade.Wait = params.Wait
	upgrade.DryRun = params.DryRun
	upgrade.ReuseValues = params.ReuseValues
	upgrade.ResetValues = params.ResetValues
	upgrade.Force = params.Force
	upgrade.Description = params.Description

	// Set timeout if specified
	if params.Timeout > 0 {
		upgrade.Timeout = time.Duration(params.Timeout) * time.Second
	}

	// Set chart path options if using repo
	if params.RepoURL != "" {
		upgrade.ChartPathOptions.RepoURL = params.RepoURL
		if params.Version != "" {
			upgrade.ChartPathOptions.Version = params.Version
		}
	}

	// Locate and load the chart
	chartPath, err := upgrade.ChartPathOptions.LocateChart(params.ChartName, s.HelmEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	// Upgrade the chart
	return upgrade.Run(params.ReleaseName, chart, params.Values)
}

// UninstallChart uninstalls a Helm chart release
func (s *Server) UninstallChart(params ChartUninstallRequest) (*release.UninstallReleaseResponse, error) {
	cfg, err := s.GetActionConfig(params.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize uninstall action
	uninstall := action.NewUninstall(cfg)
	uninstall.Wait = params.Wait
	uninstall.DryRun = params.DryRun
	uninstall.KeepHistory = params.KeepHistory
	uninstall.Description = params.Description

	// Set timeout if specified
	if params.Timeout > 0 {
		uninstall.Timeout = time.Duration(params.Timeout) * time.Second
	}

	// Uninstall the chart
	return uninstall.Run(params.ReleaseName)
}

// ListReleases returns a list of all releases
func (s *Server) ListReleases(namespace string, allNamespaces bool) ([]*release.Release, error) {
	cfg, err := s.GetActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize list action
	list := action.NewList(cfg)
	list.AllNamespaces = allNamespaces
	list.SetStateMask()

	// Return list of releases
	return list.Run()
}

// GetRelease returns details of a specific release
func (s *Server) GetRelease(name, namespace string) (*release.Release, error) {
	cfg, err := s.GetActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize get action
	get := action.NewGet(cfg)

	// Get the release
	return get.Run(name)
}

// GetReleaseHistory returns the history of a specific release
func (s *Server) GetReleaseHistory(name, namespace string) ([]*release.Release, error) {
	cfg, err := s.GetActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize history action
	history := action.NewHistory(cfg)
	history.Max = 256

	// Get release history
	return history.Run(name)
}

// GetReleaseStatus returns the status of a specific release
func (s *Server) GetReleaseStatus(name, namespace string) (*release.Release, error) {
	cfg, err := s.GetActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Initialize status action
	status := action.NewStatus(cfg)

	// Get release status
	return status.Run(name)
}

// AddRepository adds a new chart repository
func (s *Server) AddRepository(name, url string) error {
	// Create entry
	entry := repo.Entry{
		Name: name,
		URL:  url,
	}

	// Add repository to Helm configuration
	settings := cli.New()

	// Initialize chart repo with builtin providers
	providers := getter.All(settings)
	r, err := repo.NewChartRepository(&entry, providers)
	if err != nil {
		return fmt.Errorf("failed to create chart repository: %w", err)
	}

	// Set the repository cache file
	r.CachePath = settings.RepositoryCache

	// Update repository cache
	if _, err := r.DownloadIndexFile(); err != nil {
		return fmt.Errorf("failed to download repository index: %w", err)
	}

	return nil
}