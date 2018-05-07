package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"crypto/tls"
	"crypto/x509"

	"github.com/golang/glog"

	"github.com/dkoshkin/admission-webhook/pkg/admitter"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch
)

// Config contains the server (the webhook) cert, key and CA.
type Config struct {
	CertFile string
	KeyFile  string
	CAFile   string
}

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)
var decoder = codecs.UniversalDeserializer()

func addToScheme(scheme *runtime.Scheme) {
	corev1.AddToScheme(scheme)
	admissionregistrationv1beta1.AddToScheme(scheme)
}

func main() {
	// register webhook
	addToScheme(scheme)

	config := configFromFlags()
	tlsConfig := tlsConfig(config)

	glog.Fatal(server(tlsConfig))
}

func server(tlsConfig *tls.Config) error {
	http.HandleFunc("/pods", servePods)

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}
	return server.ListenAndServeTLS("", "")
}

func servePods(w http.ResponseWriter, r *http.Request) {
	a := admitter.PodHostNetworkValidator{
		UniversalDeserializerAdmitter: admitter.UniversalDeserializerAdmitter{
			Decoder: decoder,
		},
	}
	admitter.Serve(w, r, a)
}

func tlsConfig(config *Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		glog.Fatalf("error loading TLS key pair: %v:", err)
	}

	clientCACert, err := ioutil.ReadFile(config.CAFile)
	if err != nil {
		glog.Fatalf("error loading CA: %v:", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	return &tls.Config{
		ClientAuth:   tls.NoClientCert,
		Certificates: []tls.Certificate{sCert},
		RootCAs:      clientCertPool,
	}
}

func configFromFlags() *Config {
	var config Config
	flag.StringVar(&config.CAFile, "tls-ca-file", config.CAFile, "CA file used to sign the cert key pair.")
	flag.StringVar(&config.CertFile, "tls-cert-file", config.CertFile, "File containing the default x509 Certificate for HTTPS.")
	flag.StringVar(&config.KeyFile, "tls-private-key-file", config.KeyFile, "File containing the default x509 private key matching --tls-cert-file.")
	flag.Parse()
	return &config
}
