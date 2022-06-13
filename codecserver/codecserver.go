package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/temporalio/screencast-encryption-with-go/codec"
	"go.temporal.io/sdk/converter"
)

func NewCORSHTTPHandler(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Namespace")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

var originFlag string

func init() {
	flag.StringVar(&originFlag, "origin", "", "Temporal Web UI URL")
}

func main() {
	flag.Parse()

	keyID := os.Getenv("ENCRYPTION_KEY_ID")
	if keyID == "" {
		log.Fatal("Codec Server requires the ENCRYPTION_KEY_ID environment variable to be set.")
	}
	if originFlag == "" {
		log.Fatal("Codec Server requires the origin flag to enable CORS.")
	}

	handler := converter.NewPayloadCodecHTTPHandler(codec.NewEncryptionCodec(codec.CodecOptions{KeyID: keyID}))
	// Wrap the handler to add CORS support for the Temporal Web UI.
	handler = NewCORSHTTPHandler(originFlag, handler)

	http.Handle("/", handler)

	err := http.ListenAndServe(":8234", nil)
	log.Fatal(err)
}
