package k8sclient

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ClientConf struct {
	v1.TypeMeta
}

type ObjectCount struct {
	Kind string `form:"Kind" json:"Kind" yaml:"Kind" xml:"Kind" db:"Kind"`

	Namespace string `form:"namespace" json:"namespace" yaml:"namespace" xml:"namespace" db:"namespace"`

	Total int32 `form:"Kind" json:"Kind" yaml:"Kind" xml:"Kind" db:"Kind"`

	Ready int32 `form:"ready" json:"ready" yaml:"ready" xml:"ready" db:"ready"`

	Unready int32 `form:"unready" json:"unready" yaml:"unready" xml:"unready" db:"unready"`
}