package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := loadConfig()

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(conf.AWSRegionID),
	}))

	gs := mustNewGlueStore(
		conf.Tenants,
		conf.SSMPrefix,
		conf.AWSAccountID,
		conf.AWSRegionID,
		awsSession,
	)

	for _, m := range conf.Metrics {
		recorder := newMetricRecorder(gs, m)
		go recorder.startRecording()
	}

	ep := newServer(conf.ListenAddress, getRoutes())
	go func() { panicOnErr(ep.ListenAndServe(), "endpoint failed") }()
	log.WithField("core.addr", conf.ListenAddress).Info("started endpoint")

	ShutdownOnSignal(ep)
}

// ShutdownOnSignal sets up OS signal listeners and gracefully shuts down servers passed in
func ShutdownOnSignal(servers ...*http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch
	log.Infof("Server got signal %s to shutdown, doing so...", sig)

	for _, server := range servers {
		if err := server.Shutdown(context.Background()); err != nil {
			log.WithError(err).Error("unable to gracefully shutdown server")
		}
	}
}

func getRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/healthcheck", createHealthCheckEndpoint())
	mux.Handle("/metrics", createMetricsEndpoint())
	return mux
}

func createHealthCheckEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json, err := json.Marshal(map[string]string{"status": "healthy"})
		var request []byte
		_, err = r.Body.Read(request)
		if err != nil {
			log.WithError(err).Error("Failed to read body")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
		if err != nil {
			log.WithField("core.error", err).
				WithField("core.request", string(request)).
				Error("Unable to marshal healthcheck")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
}

func panicOnErr(err error, msg string, args ...interface{}) {
	if err != nil {
		log.WithError(err).Panicf(msg, args...)
	}
}
