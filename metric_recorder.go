package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type metricRecorder struct {
	gs glueStore
	m  metric
	g  *prometheus.GaugeVec
}

func (r *metricRecorder) startRecording() {
	for {
		r.record()

		time.Sleep(r.m.QueryInterval)
	}
}

func (r *metricRecorder) record() {
glueClientLoop:
	for _, gc := range r.gs.clients {
		for _, excludedDBName := range r.m.ExlcudeDBs {
			matched, err := regexp.MatchString(excludedDBName, gc.getDBName())
			if err != nil {
				log.WithError(err).Fatalf("Failed to compile regex for db name [%s]", excludedDBName)
			}
			if excludedDBName == gc.getDBName() || matched {
				continue glueClientLoop
			}
		}
		if len(r.m.IncludeDBS) > 0 {
			included := false
			for _, includedDBName := range r.m.IncludeDBS {
				matched, err := regexp.MatchString(includedDBName, gc.getDBName())
				if err != nil {
					log.WithError(err).Fatalf("Failed to compile regex for db name [%s]", includedDBName)
				}
				if includedDBName == gc.getDBName() || matched {
					included = true
					break
				}
			}
			if !included {
				continue glueClientLoop
			}
		}

		res, err := gc.runQuery(r.m.Query)
		if err != nil {
			log.WithError(err).
				Errorf(
					"Failed to run query [%s] for [%s] on [%s]",
					r.m.Query,
					gc.getTenant(),
					gc.getDBName(),
				)
			continue
		}

		filteredLabels := map[string]string{
			"tenant": gc.getTenant(),
			"db":     gc.getDBName(),
		}
	filteredLoop:
		for column, val := range res {
			for _, vColumn := range r.m.QueryValueColumns {
				if column == vColumn {
					continue filteredLoop
				}
			}

			filteredLabels[column] = val.(string)
		}

		for _, vColumn := range r.m.QueryValueColumns {
			value, ok := res[vColumn]
			if value == nil || !ok {
				log.Warnf(
					"No value found for column %s, when running query [%s] for [%s] on [%s]",
					vColumn,
					r.m.Query,
					gc.getTenant(),
					gc.getDBName(),
				)
				continue
			}

			convValue, ok := value.(float64)
			if !ok {
				v := value.(int64)
				convValue = float64(v)
			}
			filteredLabels["_column"] = vColumn

			if r.g == nil {
				filteredLabels["_column"] = vColumn
				labels := []string{}
				for label := range filteredLabels {
					labels = append(labels, label)
				}
				r.g = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: strings.ToLower(strings.Replace(r.m.Name, " ", "_", -1)),
						Help: fmt.Sprintf("Gauge for %s", r.m.Name),
					},
					labels,
				)
				mustRegisterCollector(r.g)
			}

			r.g.With(filteredLabels).Set(convValue)
		}
	}
}

func newMetricRecorder(gs glueStore, m metric) metricRecorder {
	return metricRecorder{
		gs,
		m,
		nil,
	}
}
