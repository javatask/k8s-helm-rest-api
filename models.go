package main

import (
	"time"

	"helm.sh/helm/v3/pkg/release"
)

// ChartInstallRequest defines parameters for chart installation
type ChartInstallRequest struct {
	ReleaseName     string                 `json:"releaseName"`
	ChartName       string                 `json:"chartName"`
	RepoURL         string                 `json:"repoURL,omitempty"`
	Version         string                 `json:"version,omitempty"`
	Namespace       string                 `json:"namespace,omitempty"`
	Values          map[string]interface{} `json:"values,omitempty"`
	Wait            bool                   `json:"wait,omitempty"`
	Timeout         int                    `json:"timeout,omitempty"`
	CreateNamespace bool                   `json:"createNamespace,omitempty"`
	DryRun          bool                   `json:"dryRun,omitempty"`
	ClientOnly      bool                   `json:"clientOnly,omitempty"`
	Description     string                 `json:"description,omitempty"`
}

// ChartUpgradeRequest defines parameters for chart upgrades
type ChartUpgradeRequest struct {
	ChartInstallRequest
	ReuseValues bool `json:"reuseValues,omitempty"`
	ResetValues bool `json:"resetValues,omitempty"`
	Force       bool `json:"force,omitempty"`
}

// ChartUninstallRequest defines parameters for chart uninstallation
type ChartUninstallRequest struct {
	ReleaseName  string `json:"releaseName"`
	Namespace    string `json:"namespace,omitempty"`
	KeepHistory  bool   `json:"keepHistory,omitempty"`
	Wait         bool   `json:"wait,omitempty"`
	Timeout      int    `json:"timeout,omitempty"`
	DryRun       bool   `json:"dryRun,omitempty"`
	Description  string `json:"description,omitempty"`
}

// ReleaseInfo represents summarized release information
type ReleaseInfo struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Version      int    `json:"version"`
	Status       string `json:"status"`
	LastDeployed string `json:"lastDeployed"`
	Chart        string `json:"chart"`
	AppVersion   string `json:"appVersion,omitempty"`
}

// ReleaseDetails represents detailed release information
type ReleaseDetails struct {
	Name         string                 `json:"name"`
	Namespace    string                 `json:"namespace"`
	Version      int                    `json:"version"`
	Status       string                 `json:"status"`
	Description  string                 `json:"description,omitempty"`
	FirstDeployed string                `json:"firstDeployed"`
	LastDeployed  string                `json:"lastDeployed"`
	Chart        string                 `json:"chart"`
	ChartVersion string                 `json:"chartVersion"`
	AppVersion   string                 `json:"appVersion,omitempty"`
	Values       map[string]interface{} `json:"values,omitempty"`
	Manifest     string                 `json:"manifest,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
}

// ReleaseHistoryEntry represents an entry in release history
type ReleaseHistoryEntry struct {
	Revision     int       `json:"revision"`
	Status       string    `json:"status"`
	Chart        string    `json:"chart"`
	ChartVersion string    `json:"chartVersion"`
	AppVersion   string    `json:"appVersion,omitempty"`
	Description  string    `json:"description,omitempty"`
	DeployedAt   time.Time `json:"deployedAt"`
}

// ApiResponse is a generic response object for API calls
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// FromRelease converts a Helm release to a summarized ReleaseInfo
func ReleaseInfoFromRelease(r *release.Release) ReleaseInfo {
	return ReleaseInfo{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Version:      r.Version,
		Status:       r.Info.Status.String(),
		LastDeployed: r.Info.LastDeployed.String(),
		Chart:        r.Chart.Metadata.Name,
		AppVersion:   r.Chart.Metadata.AppVersion,
	}
}

// ReleaseDetailsFromRelease converts a Helm release to detailed ReleaseDetails
func ReleaseDetailsFromRelease(r *release.Release) ReleaseDetails {
	return ReleaseDetails{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Version:      r.Version,
		Status:       r.Info.Status.String(),
		Description:  r.Info.Description,
		FirstDeployed: r.Info.FirstDeployed.String(),
		LastDeployed:  r.Info.LastDeployed.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
		Values:       r.Config,
		Manifest:     r.Manifest,
		Notes:        r.Info.Notes,
	}
}
