// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

//go:build e2e
// +build e2e

package e2e

import "embed"

var (
	//go:embed testdata/*.yaml base/*
	Manifests embed.FS

	//go:embed testdata/merge-gateways/*.yaml base/*
	MergeGatewaysManifests embed.FS
)
