// Package harconnectivity contains the HAR Connectivity network experiment.
package harconnectivity

import (
	"context"
	"errors"
	"time"

	"github.com/ooni/probe-engine/experiment"
	"github.com/ooni/probe-engine/experiment/handler"
	"github.com/ooni/probe-engine/httpx/fetch"
	"github.com/ooni/probe-engine/httpx/minihar"
	"github.com/ooni/probe-engine/model"
	"github.com/ooni/probe-engine/oohar"
	"github.com/ooni/probe-engine/session"
)

const (
	testName    = "harconnectivity"
	testVersion = "0.0.1"
)

// Config contains the experiment config.
type Config struct {
}

// TestKeys contains the experiment test keys.
type TestKeys struct {
	// Failure is the experiment result.
	Failure string `json:"failure"`

	// HAR is the HAR log for the measurement.
	HAR *oohar.HAR `json:"har"`
}

func measure(
	ctx context.Context, sess *session.Session, measurement *model.Measurement,
	callbacks handler.Callbacks, config Config,
) error {
	testkeys := new(TestKeys)
	measurement.TestKeys = testkeys
	// TODO(bassosimone): what would be a good timeout? We have a wide range
	// of network access link speeds, so it's reasonable to say that there is
	// no one-size-fits-all timeout setting here. For now, I am using 1 min
	// under the assumption that after 1 min without the measurement converging
	// a user would most likely be super pissed off. When we'll start having
	// more HAR based data, we can use data to sort this out, maybe?
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	previousrs := minihar.ContextRequestSaver(ctx)
	ctx, rs := minihar.WithRequestSaver(ctx)
	if measurement.Input == "" {
		return errors.New("harconnectivity: passed an empty input")
	}
	fetchclient := &fetch.Client{
		HTTPClient: sess.HTTPNoProxyClient,
		Logger:     sess.Logger,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
	}
	_, err := fetchclient.Fetch(ctx, measurement.Input)
	if err != nil {
		testkeys.Failure = err.Error()
	}
	testkeys.HAR = oohar.NewFromMiniHAR(
		sess.SoftwareName, sess.SoftwareVersion, rs,
	)
	if previousrs != nil {
		previousrs.RoundTrips = append(previousrs.RoundTrips, rs.RoundTrips...)
	}
	return err
}

// NewExperiment creates a new experiment.
func NewExperiment(
	sess *session.Session, config Config,
) *experiment.Experiment {
	return experiment.New(
		sess, testName, testVersion,
		func(
			ctx context.Context,
			sess *session.Session,
			measurement *model.Measurement,
			callbacks handler.Callbacks,
		) error {
			return measure(ctx, sess, measurement, callbacks, config)
		})
}