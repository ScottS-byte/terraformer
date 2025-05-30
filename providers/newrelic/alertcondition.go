// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package newrelic

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	newrelic "github.com/newrelic/newrelic-client-go/newrelic"
)

type AlertConditionGenerator struct {
	NewRelicService
}

func (g *AlertConditionGenerator) createAlertConditionResources(client *newrelic.NewRelic) error {
	alertPolicies, err := client.Alerts.ListPolicies(nil)
	if err != nil {
		return err
	}

	for _, alertPolicy := range alertPolicies {
		alertConditions, err := client.Alerts.ListConditions(alertPolicy.ID)
		if err != nil {
			return err
		}

		for _, alertCondition := range alertConditions {
			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				fmt.Sprintf("%d:%d", alertPolicy.ID, alertCondition.ID),
				fmt.Sprintf("%s-%d", normalizeResourceName(alertCondition.Name), alertCondition.ID),
				"newrelic_alert_condition",
				g.ProviderName,
				[]string{}))
		}
	}
	return nil
}

func (g *AlertConditionGenerator) createAlertNrqlConditionResources(client *newrelic.NewRelic) error {
	alertPolicies, err := client.Alerts.ListPolicies(nil)
	if err != nil {
		return err
	}

	for _, alertPolicy := range alertPolicies {
		nrqlConditions, err := client.Alerts.ListNrqlConditions(alertPolicy.ID)
		if err != nil {
			return err
		}

		for _, nrqlCondition := range nrqlConditions {
			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				fmt.Sprintf("%d:%d", alertPolicy.ID, nrqlCondition.ID),
				fmt.Sprintf("%s-%d", normalizeResourceName(nrqlCondition.Name), nrqlCondition.ID),
				"newrelic_nrql_alert_condition",
				g.ProviderName,
				[]string{}))
		}
	}
	return nil
}

func (g *AlertConditionGenerator) InitResources() error {
	client, err := g.Client()
	if err != nil {
		return err
	}

	funcs := []func(*newrelic.NewRelic) error{
		g.createAlertConditionResources,
		g.createAlertNrqlConditionResources,
	}

	for _, f := range funcs {
		err := f(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *AlertConditionGenerator) PostConvertHook() error {
	for i, resource := range g.Resources {
		if resource.InstanceInfo.Type == "newrelic_alert_condition" {
			if resource.Item["violation_close_timer"] == "0" {
				delete(g.Resources[i].Item, "violation_close_timer")
			}
		}
	}

	return nil
}
