// Copyright 2023 Shaun Mitchell

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	SolutionFileName    = "solution"
	SolutionFileKind    = "Config"
	SolutionGroup       = "arete"
	SolutionFileVersion = "v1alpha1"
	SolutionAPIVersion  = SolutionGroup + "/" + SolutionFileVersion
)

var TypeMeta = yaml.ResourceMeta{
	TypeMeta: yaml.TypeMeta{
		APIVersion: SolutionAPIVersion,
		Kind:       SolutionFileKind,
	},
}

type SolutionFile struct {
	yaml.ResourceMeta `yaml:",inline"`

	Spec *Spec `yaml:"spec,omitempty"`
	Deploy *Deploy `yaml:"deploy,omitempty"`
}

type Spec struct {
	Url string `yaml:"url,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func (s *Spec) IsEmpty() bool {
	return s == nil
}

type Requires struct {
	UseConfigConnectorSA string `yaml:"useConfigConnectorSA,omitempty"`
	Iam                  []Iam  `yaml:"iam,omitempty"`
	Services						 []Services `yaml:"services,omitempty"`
	Depends							 []Depends	`yaml:"depends,omitempty"`
}

func (r *Requires) IsEmpty() bool {
	return len(r.UseConfigConnectorSA) == 0 && len(r.Iam) == 0
}

type Services struct {
	Service string `yaml:"service,omitempty"`
	Project string `yaml:"project,omitempty"`
}

type Depends struct {
	AssetType string `yaml:"asset-type,omitempty"`
	Scope			string `yaml:"scope,omitempty"`
	Name			string `yaml:"name,omitempty"`
}

type Deploy struct {
	Stage *Stage `yaml:"stage,omitempty"`
}

func (d *Deploy) IsEmpty() bool {
	return d == nil
}

type Stage struct {
	Infra *Infra `yaml:"infra,omitempty"`
	App   *App   `yaml:"app,omitempty"`
}

func (s *Stage) IsEmpty() bool {
	return s == nil
}

type Infra struct {
	KubeContext *KubeContext `yaml:"kubeContext,omitempty"`
	Requires    Requires     `yaml:"requires,omitempty"`
}

func (i *Infra) IsEmpty() bool {
	return i == nil
}

type KubeContext struct {
	ClusterName string `yaml:"clusterName,omitempty"`
	Region      string `yaml:"region,omitempty"`
	Project     string `yaml:"project,omitempty"`
	Zone        string `yaml:"zone,omitempty"`
	InternalIP  string   `yaml:"internalIP,omitempty"`
}

func (k *KubeContext) IsEmpty() bool {
	return k == nil
}

type Iam struct {
	Role     string    `yaml:"role,omitempty"`
	Member   string    `yaml:"member,omitempty"`
	Resource *Resource `yaml:"resource,omitempty"`
}

type Resource struct {
	Level string `yaml:"level,omitempty"`
	Id    string `yaml:"id,omitempty"`
}

type App struct {
	KubeContext *KubeContext `yaml:"kubeContext,omitempty"`
}
