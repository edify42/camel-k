/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trait

import (
	"testing"

	"github.com/stretchr/testify/assert"
	serving "knative.dev/serving/pkg/apis/serving/v1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/util"
	"github.com/apache/camel-k/pkg/util/kubernetes"
)

func TestConfigureTolerationTraitMissingKey(t *testing.T) {
	environment, _ := createNominalDeploymentTolerationTest()
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)

	success, err := tolerationTrait.Configure(environment)

	assert.Equal(t, false, success)
	assert.NotNil(t, err)
}

func TestConfigureTolerationTraitMissingOperator(t *testing.T) {
	environment, _ := createNominalDeploymentTolerationTest()
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"

	success, err := tolerationTrait.Configure(environment)

	assert.Equal(t, false, success)
	assert.NotNil(t, err)
}

func TestConfigureTolerationTraitMissingValue(t *testing.T) {
	environment, _ := createNominalDeploymentTolerationTest()
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"
	tolerationTrait.Operator = "Equal"

	success, err := tolerationTrait.Configure(environment)

	assert.Equal(t, false, success)
	assert.NotNil(t, err)
}

func TestConfigureTolerationTraitMissingEffect(t *testing.T) {
	environment, _ := createNominalDeploymentTolerationTest()
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"
	tolerationTrait.Operator = "Exists"

	success, err := tolerationTrait.Configure(environment)

	assert.Equal(t, false, success)
	assert.NotNil(t, err)
}

func TestApplyPodTolerationLabelsDefault(t *testing.T) {
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"
	tolerationTrait.Operator = "Equal"
	tolerationTrait.Value = "my-value"
	tolerationTrait.Effect = "NoExecute"

	environment, deployment := createNominalDeploymentTolerationTest()
	testApplyPodTolerationLabelsDefault(t, tolerationTrait, environment, &deployment.Spec.Template.Spec.Tolerations)

	environment, knativeService := createNominalKnativeServiceTolerationTest()
	testApplyPodTolerationLabelsDefault(t, tolerationTrait, environment, &knativeService.Spec.Template.Spec.Tolerations)

	environment, cronJob := createNominalCronJobTolerationTest()
	testApplyPodTolerationLabelsDefault(t, tolerationTrait, environment, &cronJob.Spec.JobTemplate.Spec.Template.Spec.Tolerations)
}

func testApplyPodTolerationLabelsDefault(t *testing.T, trait *tolerationTrait, environment *Environment, tolerations *[]corev1.Toleration) {
	err := trait.Apply(environment)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(*tolerations))
	toleration := (*tolerations)[0]
	assert.Equal(t, "my-toleration", toleration.Key)
	assert.Equal(t, corev1.TolerationOpEqual, toleration.Operator)
	assert.Equal(t, "my-value", toleration.Value)
	assert.Equal(t, corev1.TaintEffectNoExecute, toleration.Effect)
}

func TestApplyPodTolerationLabelsTolerationSeconds(t *testing.T) {
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"
	tolerationTrait.Operator = "Exists"
	tolerationTrait.Effect = "NoExecute"
	tolerationTrait.TolerationSeconds = "300"

	environment, deployment := createNominalDeploymentTolerationTest()
	testApplyPodTolerationLabelsTolerationSeconds(t, tolerationTrait, environment, &deployment.Spec.Template.Spec.Tolerations)

	environment, knativeService := createNominalKnativeServiceTolerationTest()
	testApplyPodTolerationLabelsTolerationSeconds(t, tolerationTrait, environment, &knativeService.Spec.Template.Spec.Tolerations)

	environment, cronJob := createNominalCronJobTolerationTest()
	testApplyPodTolerationLabelsTolerationSeconds(t, tolerationTrait, environment, &cronJob.Spec.JobTemplate.Spec.Template.Spec.Tolerations)
}

func testApplyPodTolerationLabelsTolerationSeconds(t *testing.T, trait *tolerationTrait, environment *Environment, tolerations *[]corev1.Toleration) {
	err := trait.Apply(environment)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(*tolerations))
	toleration := (*tolerations)[0]
	assert.Equal(t, "my-toleration", toleration.Key)
	assert.Equal(t, corev1.TolerationOpExists, toleration.Operator)
	assert.Equal(t, corev1.TaintEffectNoExecute, toleration.Effect)
	assert.Equal(t, int64(300), *toleration.TolerationSeconds)
}

func TestApplyPodTolerationLabelsTolerationSecondsFail(t *testing.T) {
	tolerationTrait := newTolerationTrait().(*tolerationTrait)
	tolerationTrait.Enabled = util.BoolP(true)
	tolerationTrait.Key = "my-toleration"
	tolerationTrait.Operator = "Exists"
	tolerationTrait.Effect = "NoExecute"
	tolerationTrait.TolerationSeconds = "abc"

	environment, deployment := createNominalDeploymentTolerationTest()
	testApplyPodTolerationLabelsTolerationSecondsFail(t, tolerationTrait, environment, &deployment.Spec.Template.Spec.Tolerations)

	environment, knativeService := createNominalKnativeServiceTolerationTest()
	testApplyPodTolerationLabelsTolerationSecondsFail(t, tolerationTrait, environment, &knativeService.Spec.Template.Spec.Tolerations)

	environment, cronJob := createNominalCronJobTolerationTest()
	testApplyPodTolerationLabelsTolerationSecondsFail(t, tolerationTrait, environment, &cronJob.Spec.JobTemplate.Spec.Template.Spec.Tolerations)
}

func testApplyPodTolerationLabelsTolerationSecondsFail(t *testing.T, trait *tolerationTrait, environment *Environment, tolerations *[]corev1.Toleration) {
	err := trait.Apply(environment)

	assert.NotNil(t, err)
}

func createNominalDeploymentTolerationTest() (*Environment, *appsv1.Deployment) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "integration-name",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{},
		},
	}

	environment := &Environment{
		Integration: &v1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "integration-name",
			},
			Status: v1.IntegrationStatus{
				Phase: v1.IntegrationPhaseDeploying,
			},
		},
		Resources: kubernetes.NewCollection(deployment),
	}

	return environment, deployment
}

func createNominalKnativeServiceTolerationTest() (*Environment, *serving.Service) {
	knativeService := &serving.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "integration-name",
		},
		Spec: serving.ServiceSpec{},
	}

	environment := &Environment{
		Integration: &v1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "integration-name",
			},
			Status: v1.IntegrationStatus{
				Phase: v1.IntegrationPhaseDeploying,
			},
		},
		Resources: kubernetes.NewCollection(knativeService),
	}

	return environment, knativeService
}

func createNominalCronJobTolerationTest() (*Environment, *v1beta1.CronJob) {
	cronJob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "integration-name",
		},
		Spec: v1beta1.CronJobSpec{},
	}

	environment := &Environment{
		Integration: &v1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "integration-name",
			},
			Status: v1.IntegrationStatus{
				Phase: v1.IntegrationPhaseDeploying,
			},
		},
		Resources: kubernetes.NewCollection(cronJob),
	}

	return environment, cronJob
}
