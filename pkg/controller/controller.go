/*
Copyright 2016 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/fission/fission/pkg/crd"
)

func Start(ctx context.Context, logger *zap.Logger, port int, unitTestFlag bool) {
	cLogger := logger.Named("controller")

	fc, kc, apiExtClient, _, err := crd.MakeFissionClient()
	if err != nil {
		cLogger.Fatal("failed to connect to k8s API", zap.Error(err))
	}

	err = crd.EnsureFissionCRDs(ctx, cLogger, apiExtClient)
	if err != nil {
		cLogger.Fatal("failed to find fission CRDs", zap.Error(err))
	}

	err = crd.WaitForCRDs(ctx, logger, fc)
	if err != nil {
		cLogger.Fatal("error waiting for CRDs", zap.Error(err))
	}

	featureStatus, err := ConfigureFeatures(ctx, cLogger, unitTestFlag, fc, kc)
	if err != nil {
		cLogger.Error("error configuring features - proceeding without optional features", zap.Error(err))
	}

	api, err := MakeAPI(cLogger, featureStatus)
	if err != nil {
		cLogger.Fatal("failed to start controller", zap.Error(err))
	}
	api.Serve(ctx, port)
}
