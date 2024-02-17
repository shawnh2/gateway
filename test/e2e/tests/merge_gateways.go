// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

//go:build e2e
// +build e2e

package tests

import (
	"net"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	"sigs.k8s.io/gateway-api/conformance/utils/kubernetes"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests, MergeGatewaysTest)
}

var MergeGatewaysTest = suite.ConformanceTest{
	ShortName:   "MergeGateways",
	Description: "Merge gateways on to a single EnvoyProxy",
	Manifests:   []string{"testdata/merge-gateways.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ns := "gateway-conformance-infra"

		route1NN := types.NamespacedName{Name: "merged-gateway-route-1", Namespace: ns}
		gw1NN := types.NamespacedName{Name: "merged-gateway-1", Namespace: ns}
		gw1HostPort := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw1NN), route1NN)
		gw1Addr, _, err := net.SplitHostPort(gw1HostPort)
		if err != nil {
			t.Errorf("failed to split the hostport of %s: %s, reason: %v", gw1NN.String(), gw1HostPort, err)
		}

		route2NN := types.NamespacedName{Name: "merged-gateway-route-2", Namespace: ns}
		gw2NN := types.NamespacedName{Name: "merged-gateway-2", Namespace: ns}
		gw2HostPort := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw2NN), route2NN)
		gw2Addr, _, err := net.SplitHostPort(gw2HostPort)
		if err != nil {
			t.Errorf("failed to split the hostport of %s: %s, reason: %v", gw2NN.String(), gw2HostPort, err)
		}

		route3NN := types.NamespacedName{Name: "merged-gateway-route-3", Namespace: ns}
		gw3NN := types.NamespacedName{Name: "merged-gateway-3", Namespace: ns}
		gw3HostPort := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw3NN), route3NN)
		gw3Addr, _, err := net.SplitHostPort(gw3HostPort)
		if err != nil {
			t.Errorf("failed to split the hostport of %s: %s, reason: %v", gw3NN.String(), gw3HostPort, err)
		}

		t.Run("merged three gateways under the same namespace with http routes", func(t *testing.T) {
			if gw1Addr != gw2Addr {
				t.Errorf("failed to merge gateways: inconsistent gateway address %s and %s for %s and %s", gw1Addr, gw2Addr, gw1NN.String(), gw2NN.String())
				t.FailNow()
			}
			if gw2Addr != gw3Addr {
				t.Errorf("failed to merge gateways: inconsistent gateway address %s and %s for %s and %s", gw2Addr, gw3Addr, gw2NN.String(), gw3NN.String())
				t.FailNow()
			}

			// Three gateways should have the same address.
			gwAddr := gw1Addr

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge1"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns,
				Backend:   "infra-backend-v1",
			})

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge2"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns,
				Backend:   "infra-backend-v2",
			})

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge3"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns,
				Backend:   "infra-backend-v3",
			})
		})
	},
}
