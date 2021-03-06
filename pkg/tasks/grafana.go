// Copyright 2018 The Cluster Monitoring Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type GrafanaTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewGrafanaTask(client *client.Client, factory *manifests.Factory) *GrafanaTask {
	return &GrafanaTask{
		client:  client,
		factory: factory,
	}
}

func (t *GrafanaTask) Run() error {
	cr, err := t.factory.GrafanaClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana ClusterRole failed")
	}

	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana ClusterRole failed")
	}

	crb, err := t.factory.GrafanaClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana ClusterRoleBinding failed")
	}

	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana ClusterRoleBinding failed")
	}

	r, err := t.factory.GrafanaRoute()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Route failed")
	}

	err = t.client.CreateRouteIfNotExists(r)
	if err != nil {
		return errors.Wrap(err, "creating Grafana Route failed")
	}

	_, err = t.client.WaitForRouteReady(r)
	if err != nil {
		return errors.Wrap(err, "waiting for Grafana Route to become ready failed")
	}

	ps, err := t.factory.GrafanaProxySecret()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana proxy Secret failed")
	}

	err = t.client.CreateIfNotExistSecret(ps)
	if err != nil {
		return errors.Wrap(err, "creating Grafana proxy Secret failed")
	}

	cmc, err := t.factory.GrafanaConfigConfigMap()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Config ConfigMap failed")
	}

	err = t.client.CreateOrUpdateConfigMap(cmc)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana Config ConfigMap failed")
	}

	cmds, err := t.factory.GrafanaDatasources()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Datasources ConfigMap failed")
	}

	err = t.client.CreateIfNotExistConfigMap(cmds)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana Datasources ConfigMap failed")
	}

	cmdds, err := t.factory.GrafanaDashboardDefinitions()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Dashboard Definitions ConfigMaps failed")
	}

	err = t.client.CreateOrUpdateConfigMapList(cmdds)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana Dashboard Definitions ConfigMaps failed")
	}

	cmdbs, err := t.factory.GrafanaDashboardSources()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Dashboard Sources ConfigMap failed")
	}

	err = t.client.CreateOrUpdateConfigMap(cmdbs)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana Dashboard Sources ConfigMap failed")
	}

	sa, err := t.factory.GrafanaServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana ServiceAccount failed")
	}

	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana ServiceAccount failed")
	}

	svc, err := t.factory.GrafanaService()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Grafana Service failed")
	}

	d, err := t.factory.GrafanaDeployment()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Deployment failed")
	}

	err = t.client.CreateOrUpdateDeployment(d)
	return errors.Wrap(err, "reconciling Grafana Deployment failed")
}
