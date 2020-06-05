package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto will perform a DeepCopy into the provided AWSProviderSpec
func (in *AWSProviderSpec) DeepCopyInto(out *AWSProviderSpec) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.StatementEntries != nil {
		in, out := &in.StatementEntries, &out.StatementEntries
		*out = make([]StatementEntry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy will DeepCopy and return a pointer to a
// new AWSProviderSpec
func (in *AWSProviderSpec) DeepCopy() *AWSProviderSpec {
	if in == nil {
		return nil
	}
	out := new(AWSProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject will return a DeepCopied AWSProviderSpec
// as a runtime.Object
func (in *AWSProviderSpec) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func deepCopyIAMPolicyCondition(ipc IAMPolicyCondition) IAMPolicyCondition {
	cp := make(IAMPolicyCondition)
	for key, val := range ipc {
		if val != nil {
			cp[key] = make(IAMPolicyConditionKeyValue)
			for subKey, subVal := range val {
				cp[key][subKey] = subVal
			}
		}
	}

	return cp
}

// DeepCopyInto will perform a DeepCopy into the provided StatementEntry
func (in *StatementEntry) DeepCopyInto(out *StatementEntry) {
	*out = *in
	if in.Action != nil {
		in, out := &in.Action, &out.Action
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PolicyCondition != nil {
		out.PolicyCondition = deepCopyIAMPolicyCondition(in.PolicyCondition)
	}

	return
}

// DeepCopy will DeepCopy and return a pointer to a
// new StatementEntry
func (in *StatementEntry) DeepCopy() *StatementEntry {
	if in == nil {
		return nil
	}
	out := new(StatementEntry)
	in.DeepCopyInto(out)
	return out
}
