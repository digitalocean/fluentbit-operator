package filter

import (
	"testing"

	. "github.com/onsi/gomega"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"
)

func TestOutput_RewriteTag_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace", nil)

	rt := RewriteTag{
		EmitterName: "my_emitter",
		Rules: []string{
			"$tool ^(fluent)$  from.$TAG.new.$tool.$sub['s1']['s2'].out false",
			"$log ^.*$  logtail.logs.new true",
		},
		EmitterStorageType:       "memory",
		EmitterMemoryBufferLimit: "10M",
	}

	expected := params.NewKVs()
	expected.Insert("Rule", "$tool ^(fluent)$  from.$TAG.new.$tool.$sub['s1']['s2'].out false")
	expected.Insert("Rule", "$log ^.*$  logtail.logs.new true")
	expected.Insert("Emitter_Name", "my_emitter")
	expected.Insert("Emitter_Storage.type", "memory")
	expected.Insert("Emitter_Mem_Buf_Limit", "10M")

	kvs, err := rt.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
