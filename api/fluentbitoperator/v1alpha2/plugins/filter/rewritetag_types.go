package filter

import (
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// RewriteTag define a `rewrite_tag` filter, allows to re-emit a record under a new Tag.
// Once a record has been re-emitted, the original record can be preserved or discarded.
type RewriteTag struct {
	// Defines the matching criteria and the format of the Tag for the matching record.
	// The Rule format have four components: KEY REGEX NEW_TAG KEEP.
	Rules []string `json:"rules,omitempty"`
	// When the filter emits a record under the new Tag, there is an internal emitter
	// plugin that takes care of the job. Since this emitter expose metrics as any other
	// component of the pipeline, you can use this property to configure an optional name for it.
	EmitterName string `json:"emitterName,omitempty"`
	// Define a buffering mechanism for the new records created.
	// Note these records are part of the emitter plugin. This option support the values
	// "memory" (default) or "filesystem".
	EmitterStorageType string `json:"emitterStorageType,omitempty"`
	// Set a limit on the amount of memory the tag rewrite emitter can consume
	// if the outputs provide backpressure.
	EmitterMemoryBufferLimit string `json:"emitterMemoryBufferLimit,omitempty"`
}

func (_ *RewriteTag) Name() string {
	return "rewrite_tag"
}

func (r *RewriteTag) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	for _, rule := range r.Rules {
		kvs.Insert("Rule", rule)
	}
	if r.EmitterName != "" {
		kvs.Insert("Emitter_Name", r.EmitterName)
	}
	if r.EmitterStorageType != "" {
		kvs.Insert("Emitter_Storage.type", r.EmitterStorageType)
	}
	if r.EmitterMemoryBufferLimit != "" {
		kvs.Insert("Emitter_Mem_Buf_Limit", r.EmitterMemoryBufferLimit)
	}

	return kvs, nil
}
