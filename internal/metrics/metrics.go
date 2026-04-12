// Package metrics
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LinkGenerateTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_link_generate_total",
			Help: "Generate links request by result: created|existing|error",
		},
		[]string{"result"},
	)

	LinkInfoTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_link_info_total",
			Help: "Get link info requests by result: ok|not_found|expired|no_attempts|crm_error|error",
		},
		[]string{"result"},
	)

	DocPermissionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_doc_permission_total",
			Help: "Doc permission(OTP) requests by result: granted|not_found|wrong_otp|no_attempts|crm_error|error",
		},
		[]string{"result"},
	)

	DocDownloadTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_doc_download_total",
			Help: "Doc download requests by result: success|no_doc|wrong_token|wrong_status|crm_error|error",
		},
		[]string{"result"},
	)

	DocSignTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_doc_sign_total",
			Help: "Doc sign requests by result: success|wrong_status|no_doc|no_pep|wrong_token|crm_error|error",
		},
		[]string{"result"},
	)

	AcceptConditionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_accept_conditions_total",
			Help: "Get accept conditions by result: success|not_found|wrong_token|pep_exists|wrong_status|crm_error|error",
		},
		[]string{"result"})

	TotalLinks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "paperless_total_links",
			Help: "Current number of generated links",
		},
	)

	ActiveLinks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "paperless_active_links",
			Help: "Current number of active links",
		},
	)

	CRMRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_crm_requests_total",
			Help: "CRM requests by operation and result",
		},
		[]string{"operation", "result"},
	)

	CRMRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "paperless_crm_request_duration_seconds",
			Help:    "CRM request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "result"},
	)

	CRMRetriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "paperless_crm_retries_total",
			Help: "CRM retries by operation and reason",
		},
		[]string{"operation", "reason"},
	)

	LatencyAPI = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "paperless_api_latency_seconds",
			Help:    "API endpoint latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status_class"},
	)
)
