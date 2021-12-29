package metrics

import "github.com/ory/ladon"

type PrometheusMetrics struct{}

func (mtr *PrometheusMetrics) RequestDeniedBy(r ladon.Request, p ladon.Policy)           {}
func (mtr *PrometheusMetrics) RequestAllowedBy(r ladon.Request, policies ladon.Policies) {}
func (mtr *PrometheusMetrics) RequestNoMatch(r ladon.Request)                            {}
func (mtr *PrometheusMetrics) RequestProcessingError(r ladon.Request, err error)         {}
