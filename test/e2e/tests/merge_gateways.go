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
	MergeGatewaysTests = append(MergeGatewaysTests, MergeGatewaysTest)
}

var MergeGatewaysTest = suite.ConformanceTest{
	ShortName:   "BasicMergeGateways",
	Description: "Basic test for MergeGateways feature",
	Manifests:   []string{"testdata/basic-merge-gateways.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ns := "gateway-conformance-app-backend"

		route1NN := types.NamespacedName{Name: "merged-gateway-route-1", Namespace: ns}
		gw1NN := types.NamespacedName{Name: "merged-gateway-1", Namespace: ns}
		gw1HostPort := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw1NN), route1NN)

		route2NN := types.NamespacedName{Name: "merged-gateway-route-2", Namespace: ns}
		gw2NN := types.NamespacedName{Name: "merged-gateway-2", Namespace: ns}
		gw2HostPort := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gw2NN), route2NN)

		t.Run("merged gateways under the same namespace with http routes", func(t *testing.T) {
			gw1Addr, _, err := net.SplitHostPort(gw1HostPort)
			if err != nil {
				t.Errorf("failed to split hostport %s of gateway %s: %v", gw1HostPort, gw1NN.Name, err)
			}

			gw2Addr, _, err := net.SplitHostPort(gw2HostPort)
			if err != nil {
				t.Errorf("failed to split hostport %s of gateway %s: %v", gw2HostPort, gw2NN.Name, err)
			}

			if gw1Addr != gw2Addr {
				t.Errorf("inconsistent gateway address %s and %s for %s and %s", gw1Addr, gw2Addr, gw1NN.String(), gw2NN.String())
				t.FailNow()
			}

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gw1HostPort, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge1", Host: "www.example1.com"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns,
				Backend:   "app-backend-v1",
			})

			http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gw2HostPort, http.ExpectedResponse{
				Request:   http.Request{Path: "/merge2", Host: "www.example2.com"},
				Response:  http.Response{StatusCode: 200},
				Namespace: ns,
				Backend:   "app-backend-v2",
			})
		})
	},
}
