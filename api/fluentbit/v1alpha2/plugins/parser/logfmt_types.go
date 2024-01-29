package parser

import (
	"kubesphere.io/fluentbit-operator/api/fluentbit/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The logfmt parser allows to parse the logfmt format described in https://brandur.org/logfmt . <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/parsers/logfmt**
type Logfmt struct{}

func (_ *Logfmt) Name() string {
	return "logfmt"
}

func (_ *Logfmt) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
