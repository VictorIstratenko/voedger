/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 */

package istructsmem

import (
	"sync"

	"github.com/voedger/voedger/pkg/irates"
	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istructs"
	payloads "github.com/voedger/voedger/pkg/itokens-payloads"
)

// Provide: constructs new application structures provider
func Provide(appConfigs AppConfigsType, bucketsFactory irates.BucketsFactoryType, appTokensFactory payloads.IAppTokensFactory,
	storageProvider istorage.IAppStorageProvider) (provider istructs.IAppStructsProvider) {
	return &appStructsProviderType{
		locker:           sync.RWMutex{},
		configs:          appConfigs,
		structures:       make(map[istructs.AppQName]*appStructsType),
		bucketsFacotry:   bucketsFactory,
		appTokensFactory: appTokensFactory,
		storageProvider:  storageProvider,
	}
}
