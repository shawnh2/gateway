//go:build conformance
// +build conformance

// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package conformance

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/flags"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func TestGatewayAPIConformance(t *testing.T) {
	flag.Parse()

	cfg, err := config.GetConfig()
	require.NoError(t, err)

	client, err := client.New(cfg, client.Options{})
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	require.NoError(t, v1alpha2.AddToScheme(client.Scheme()))
	require.NoError(t, v1beta1.AddToScheme(client.Scheme()))

	cSuite := suite.New(suite.Options{
		Client:               client,
		GatewayClassName:     *flags.GatewayClassName,
		Debug:                *flags.ShowDebug,
		Clientset:            clientset,
		CleanupBaseResources: *flags.CleanupBaseResources,
		SupportedFeatures:    suite.AllFeatures,
		ExemptFeatures:       suite.MeshCoreFeatures,
	})
	cSuite.Setup(t)
	cSuite.Run(t, tests.ConformanceTests)

	// --- FAIL: TestGatewayAPIConformance/HTTPRouteRedirectPortAndScheme (1.11s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/0_request_to_'/scheme-nil-and-port-nil'_should_receive_a_302 (0.01s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/http-listener-on-8080/0_request_to_'/scheme-nil-and-port-nil'_should_receive_a_302 (30.00s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/5_request_to_'example.org/scheme-http-and-port-8080'_should_receive_a_302 (31.51s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/4_request_to_'example.org/scheme-http-and-port-80'_should_receive_a_302 (30.00s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/3_request_to_'example.org/scheme-http-and-port-nil'_should_receive_a_302 (30.72s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/2_request_to_'example.org/scheme-nil-and-port-8443'_should_receive_a_302 (30.00s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/1_request_to_'example.org/scheme-nil-and-port-443'_should_receive_a_302 (30.72s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-8080/2_request_to_'/scheme-https-and-port-nil'_should_receive_a_302 (0.00s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-8080/1_request_to_'/scheme-nil-and-port-80'_should_receive_a_302 (0.00s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/2_request_to_'/scheme-nil-and-port-8080'_should_receive_a_302 (0.00s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/5_request_to_'/scheme-https-and-port-8443'_should_receive_a_302 (0.00s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/4_request_to_'/scheme-https-and-port-443'_should_receive_a_302 (0.01s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/3_request_to_'/scheme-https-and-port-nil'_should_receive_a_302 (0.00s)
	//    --- PASS: HTTPRouteRedirectPortAndScheme/http-listener-on-80/1_request_to_'/scheme-nil-and-port-80'_should_receive_a_302 (0.00s)
	//    --- FAIL: HTTPRouteRedirectPortAndScheme/https-listener-on-443/0_request_to_'example.org/scheme-nil-and-port-nil'_should_receive_a_302 (30.59s)

}
