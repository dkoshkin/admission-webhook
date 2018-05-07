package admitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Serve(w http.ResponseWriter, r *http.Request, admitter Admitter) {
	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err != nil {
			glog.Errorf("error reading body: %v", err)
			return
		} else {
			body = data
		}
	}

	glog.V(2).Info(fmt.Sprintf("handling request: %v", body))
	var reviewResponse *v1beta1.AdmissionResponse

	ar, err := admitter.Decode(body)
	if err != nil {
		glog.Errorf("error decoding body: %v", err)
		reviewResponse = toAdmissionResponse(err)
	} else {
		reviewResponse, err = admitter.Admit(ar)
		if err != nil {
			glog.Error(err)
			reviewResponse = toAdmissionResponse(err)
		}
	}
	glog.V(2).Info(fmt.Sprintf("sending response: %v", reviewResponse))

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = ar.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	ar.Request.Object = runtime.RawExtension{}
	ar.Request.OldObject = runtime.RawExtension{}

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Errorf("error marshalling response: %v", err)
		return
	}
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("error writing response: %v", err)
		return
	}
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
