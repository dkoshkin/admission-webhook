package admitter

import (
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

type Admitter interface {
	Decode(body []byte) (*v1beta1.AdmissionReview, error)
	Admit(ar *v1beta1.AdmissionReview) (*v1beta1.AdmissionResponse, error)
}

type UniversalDeserializerAdmitter struct {
	Decoder runtime.Decoder
}

func (a UniversalDeserializerAdmitter) Decode(body []byte) (*v1beta1.AdmissionReview, error) {
	ar := v1beta1.AdmissionReview{}
	if _, _, err := a.Decoder.Decode(body, nil, &ar); err != nil {
		return nil, err
	}
	return &ar, nil

}
