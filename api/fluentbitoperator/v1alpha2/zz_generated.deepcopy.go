//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/filter"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/input"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/output"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/parser"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Decorder) DeepCopyInto(out *Decorder) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decorder.
func (in *Decorder) DeepCopy() *Decorder {
	if in == nil {
		return nil
	}
	out := new(Decorder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filter) DeepCopyInto(out *Filter) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filter.
func (in *Filter) DeepCopy() *Filter {
	if in == nil {
		return nil
	}
	out := new(Filter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Filter) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FilterItem) DeepCopyInto(out *FilterItem) {
	*out = *in
	if in.Grep != nil {
		in, out := &in.Grep, &out.Grep
		*out = new(filter.Grep)
		**out = **in
	}
	if in.RecordModifier != nil {
		in, out := &in.RecordModifier, &out.RecordModifier
		*out = new(filter.RecordModifier)
		(*in).DeepCopyInto(*out)
	}
	if in.Kubernetes != nil {
		in, out := &in.Kubernetes, &out.Kubernetes
		*out = new(filter.Kubernetes)
		(*in).DeepCopyInto(*out)
	}
	if in.Modify != nil {
		in, out := &in.Modify, &out.Modify
		*out = new(filter.Modify)
		(*in).DeepCopyInto(*out)
	}
	if in.Nest != nil {
		in, out := &in.Nest, &out.Nest
		*out = new(filter.Nest)
		(*in).DeepCopyInto(*out)
	}
	if in.Parser != nil {
		in, out := &in.Parser, &out.Parser
		*out = new(filter.Parser)
		(*in).DeepCopyInto(*out)
	}
	if in.Lua != nil {
		in, out := &in.Lua, &out.Lua
		*out = new(filter.Lua)
		(*in).DeepCopyInto(*out)
	}
	if in.Throttle != nil {
		in, out := &in.Throttle, &out.Throttle
		*out = new(filter.Throttle)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FilterItem.
func (in *FilterItem) DeepCopy() *FilterItem {
	if in == nil {
		return nil
	}
	out := new(FilterItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FilterList) DeepCopyInto(out *FilterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Filter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FilterList.
func (in *FilterList) DeepCopy() *FilterList {
	if in == nil {
		return nil
	}
	out := new(FilterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FilterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FilterSpec) DeepCopyInto(out *FilterSpec) {
	*out = *in
	if in.FilterItems != nil {
		in, out := &in.FilterItems, &out.FilterItems
		*out = make([]FilterItem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FilterSpec.
func (in *FilterSpec) DeepCopy() *FilterSpec {
	if in == nil {
		return nil
	}
	out := new(FilterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBit) DeepCopyInto(out *FluentBit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBit.
func (in *FluentBit) DeepCopy() *FluentBit {
	if in == nil {
		return nil
	}
	out := new(FluentBit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentBit) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitConfig) DeepCopyInto(out *FluentBitConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitConfig.
func (in *FluentBitConfig) DeepCopy() *FluentBitConfig {
	if in == nil {
		return nil
	}
	out := new(FluentBitConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentBitConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitConfigList) DeepCopyInto(out *FluentBitConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FluentBitConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitConfigList.
func (in *FluentBitConfigList) DeepCopy() *FluentBitConfigList {
	if in == nil {
		return nil
	}
	out := new(FluentBitConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentBitConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitConfigSpec) DeepCopyInto(out *FluentBitConfigSpec) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(Service)
		(*in).DeepCopyInto(*out)
	}
	in.InputSelector.DeepCopyInto(&out.InputSelector)
	in.FilterSelector.DeepCopyInto(&out.FilterSelector)
	in.OutputSelector.DeepCopyInto(&out.OutputSelector)
	in.ParserSelector.DeepCopyInto(&out.ParserSelector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitConfigSpec.
func (in *FluentBitConfigSpec) DeepCopy() *FluentBitConfigSpec {
	if in == nil {
		return nil
	}
	out := new(FluentBitConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitList) DeepCopyInto(out *FluentBitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FluentBit, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitList.
func (in *FluentBitList) DeepCopy() *FluentBitList {
	if in == nil {
		return nil
	}
	out := new(FluentBitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FluentBitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitSpec) DeepCopyInto(out *FluentBitSpec) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	in.PositionDB.DeepCopyInto(&out.PositionDB)
	in.Resources.DeepCopyInto(&out.Resources)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitSpec.
func (in *FluentBitSpec) DeepCopy() *FluentBitSpec {
	if in == nil {
		return nil
	}
	out := new(FluentBitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentBitStatus) DeepCopyInto(out *FluentBitStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentBitStatus.
func (in *FluentBitStatus) DeepCopy() *FluentBitStatus {
	if in == nil {
		return nil
	}
	out := new(FluentBitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Input) DeepCopyInto(out *Input) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Input.
func (in *Input) DeepCopy() *Input {
	if in == nil {
		return nil
	}
	out := new(Input)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Input) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InputList) DeepCopyInto(out *InputList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Input, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InputList.
func (in *InputList) DeepCopy() *InputList {
	if in == nil {
		return nil
	}
	out := new(InputList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InputList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InputSpec) DeepCopyInto(out *InputSpec) {
	*out = *in
	if in.Dummy != nil {
		in, out := &in.Dummy, &out.Dummy
		*out = new(input.Dummy)
		(*in).DeepCopyInto(*out)
	}
	if in.Tail != nil {
		in, out := &in.Tail, &out.Tail
		*out = new(input.Tail)
		(*in).DeepCopyInto(*out)
	}
	if in.Systemd != nil {
		in, out := &in.Systemd, &out.Systemd
		*out = new(input.Systemd)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InputSpec.
func (in *InputSpec) DeepCopy() *InputSpec {
	if in == nil {
		return nil
	}
	out := new(InputSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Output) DeepCopyInto(out *Output) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Output.
func (in *Output) DeepCopy() *Output {
	if in == nil {
		return nil
	}
	out := new(Output)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Output) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputList) DeepCopyInto(out *OutputList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputList.
func (in *OutputList) DeepCopy() *OutputList {
	if in == nil {
		return nil
	}
	out := new(OutputList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OutputList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputSpec) DeepCopyInto(out *OutputSpec) {
	*out = *in
	if in.Elasticsearch != nil {
		in, out := &in.Elasticsearch, &out.Elasticsearch
		*out = new(output.Elasticsearch)
		(*in).DeepCopyInto(*out)
	}
	if in.File != nil {
		in, out := &in.File, &out.File
		*out = new(output.File)
		**out = **in
	}
	if in.Forward != nil {
		in, out := &in.Forward, &out.Forward
		*out = new(output.Forward)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(output.HTTP)
		(*in).DeepCopyInto(*out)
	}
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		*out = new(output.Kafka)
		(*in).DeepCopyInto(*out)
	}
	if in.Null != nil {
		in, out := &in.Null, &out.Null
		*out = new(output.Null)
		**out = **in
	}
	if in.Stdout != nil {
		in, out := &in.Stdout, &out.Stdout
		*out = new(output.Stdout)
		**out = **in
	}
	if in.TCP != nil {
		in, out := &in.TCP, &out.TCP
		*out = new(output.TCP)
		(*in).DeepCopyInto(*out)
	}
	if in.Loki != nil {
		in, out := &in.Loki, &out.Loki
		*out = new(output.Loki)
		(*in).DeepCopyInto(*out)
	}
	if in.Syslog != nil {
		in, out := &in.Syslog, &out.Syslog
		*out = new(output.Syslog)
		(*in).DeepCopyInto(*out)
	}
	if in.DataDog != nil {
		in, out := &in.DataDog, &out.DataDog
		*out = new(output.DataDog)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputSpec.
func (in *OutputSpec) DeepCopy() *OutputSpec {
	if in == nil {
		return nil
	}
	out := new(OutputSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parser) DeepCopyInto(out *Parser) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parser.
func (in *Parser) DeepCopy() *Parser {
	if in == nil {
		return nil
	}
	out := new(Parser)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Parser) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParserList) DeepCopyInto(out *ParserList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Parser, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParserList.
func (in *ParserList) DeepCopy() *ParserList {
	if in == nil {
		return nil
	}
	out := new(ParserList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ParserList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParserSpec) DeepCopyInto(out *ParserSpec) {
	*out = *in
	if in.JSON != nil {
		in, out := &in.JSON, &out.JSON
		*out = new(parser.JSON)
		(*in).DeepCopyInto(*out)
	}
	if in.Regex != nil {
		in, out := &in.Regex, &out.Regex
		*out = new(parser.Regex)
		(*in).DeepCopyInto(*out)
	}
	if in.LTSV != nil {
		in, out := &in.LTSV, &out.LTSV
		*out = new(parser.LSTV)
		(*in).DeepCopyInto(*out)
	}
	if in.Logfmt != nil {
		in, out := &in.Logfmt, &out.Logfmt
		*out = new(parser.Logfmt)
		**out = **in
	}
	if in.Decoders != nil {
		in, out := &in.Decoders, &out.Decoders
		*out = make([]Decorder, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParserSpec.
func (in *ParserSpec) DeepCopy() *ParserSpec {
	if in == nil {
		return nil
	}
	out := new(ParserSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	if in.Daemon != nil {
		in, out := &in.Daemon, &out.Daemon
		*out = new(bool)
		**out = **in
	}
	if in.FlushSeconds != nil {
		in, out := &in.FlushSeconds, &out.FlushSeconds
		*out = new(int64)
		**out = **in
	}
	if in.GraceSeconds != nil {
		in, out := &in.GraceSeconds, &out.GraceSeconds
		*out = new(int64)
		**out = **in
	}
	if in.HttpPort != nil {
		in, out := &in.HttpPort, &out.HttpPort
		*out = new(int32)
		**out = **in
	}
	if in.HttpServer != nil {
		in, out := &in.HttpServer, &out.HttpServer
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Service.
func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}
