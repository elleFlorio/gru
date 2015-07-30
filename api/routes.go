package api

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	//SERVICE
	Route{
		"InfoServices",
		"GET",
		"/gru/v1/services",
		GetInfoServices,
	},

	// STATS
	Route{
		"StatsServices",
		"GET",
		"/gru/v1/stats/services",
		GetStatsServices,
	},

	Route{
		"StatsService",
		"GET",
		"/gru/v1/stats/services/{name}",
		GetStatsService,
	},

	Route{
		"StatsInstances",
		"GET",
		"/gru/v1/stats/instances",
		GetStatsInstances,
	},

	Route{
		"StatsInstance",
		"GET",
		"/gru/v1/stats/instances/{id}",
		GetStatsInstance,
	},

	Route{
		"StatsSystem",
		"GET",
		"/gru/v1/stats/system",
		GetStatsSystem,
	},

	//ANALYTICS
	Route{
		"AnalyticsServices",
		"GET",
		"/gru/v1/analytics/services",
		GetAnalyticsServices,
	},

	Route{
		"AnalyticsService",
		"GET",
		"/gru/v1/analytics/services/{name}",
		GetAnalyticsService,
	},

	Route{
		"AnalyticsInstances",
		"GET",
		"/gru/v1/analytics/instances",
		GetAnalyticsInstances,
	},

	Route{
		"AnalyticsInstance",
		"GET",
		"/gru/v1/analytics/instances/{id}",
		GetAnalyticsInstance,
	},

	Route{
		"AnalyticsSystem",
		"GET",
		"/gru/v1/analytics/system",
		GetAnalyticsSystem,
	},

	//NODE
	Route{
		"InfoNode",
		"GET",
		"/gru/v1/node",
		GetInfoNode,
	},

	//POLICY
	Route{
		"InfoPolicies",
		"GET",
		"/gru/v1/policies",
		GetInfoPolicies,
	},

	Route{
		"InfoPoliciesByType",
		"GET",
		"/gru/v1/policies/{type}",
		GetInfoPoliciesByType,
	},

	//ACTION
	Route{
		"InfoActions",
		"GET",
		"/gru/v1/actions",
		GetInfoActions,
	},
}
