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
		"StatsNode",
		"GET",
		"/gru/v1/stats",
		GetStatsNode,
	},

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
		"AnalyticsNode",
		"GET",
		"/gru/v1/analytics",
		GetAnalyticsNode,
	},

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
		"AnalyticsSystem",
		"GET",
		"/gru/v1/analytics/system",
		GetAnalyticsSystem,
	},

	Route{
		"AnalyticsCluster",
		"GET",
		"/gru/v1/analytics/cluster",
		GetAnalyticsCluster,
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

	//ACTION
	Route{
		"InfoActions",
		"GET",
		"/gru/v1/actions",
		GetInfoActions,
	},

	//SHARED
	Route{
		"SharedData",
		"GET",
		"/gru/v1/shared",
		GetSharedData,
	},

	//COMMANDS
	Route{
		"ExecCommands",
		"POST",
		"/gru/v1/commands",
		PostCommand,
	},
}
