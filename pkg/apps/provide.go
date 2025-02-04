/*
 * Copyright (c) 2022-present unTill Pro, Ltd.
 */

package apps

import (
	sysmonitor "github.com/voedger/voedger/pkg/apps/sys.monitor"
	"github.com/voedger/voedger/pkg/ihttpctl"
)

func ProvideStaticEmbeddedResources() []ihttpctl.StaticResourcesType {
	return []ihttpctl.StaticResourcesType{
		sysmonitor.Provide(),
	}
}
