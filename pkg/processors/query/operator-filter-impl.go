/*
 * Copyright (c) 2021-present unTill Pro, Ltd.
 */

package queryprocessor

import (
	"context"
	"time"

	pipeline "github.com/heeus/core-pipeline"
	"github.com/untillpro/voedger/pkg/istructs"
	coreutils "github.com/untillpro/voedger/pkg/utils"
)

type FilterOperator struct {
	pipeline.AsyncNOOP
	filters    []IFilter
	rootSchema coreutils.SchemaFields
	metrics    IMetrics
}

func (o FilterOperator) DoAsync(ctx context.Context, work pipeline.IWorkpiece) (outWork pipeline.IWorkpiece, err error) {
	begin := time.Now()
	defer func() {
		o.metrics.Increase(execFilterSeconds, time.Since(begin).Seconds())
	}()
	outputRow := work.(IWorkpiece).OutputRow().Value(rootDocument).([]IOutputRow)[0]
	mergedSchema := make(map[string]istructs.DataKindType)
	for k, v := range o.rootSchema {
		mergedSchema[k] = v
	}
	for k, v := range work.(IWorkpiece).EnrichedRootSchema() {
		mergedSchema[k] = v
	}
	for _, filter := range o.filters {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		match, err := filter.IsMatch(mergedSchema, outputRow)
		if err != nil {
			return nil, err
		}
		if !match {
			work.Release()
			return nil, nil
		}
	}
	return work, nil
}
