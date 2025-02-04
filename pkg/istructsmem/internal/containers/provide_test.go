/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package containers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istorageimpl"
	"github.com/voedger/voedger/pkg/istructs"
	"github.com/voedger/voedger/pkg/istructsmem/internal/vers"
)

func Test_BasicUsage(t *testing.T) {
	sp := istorageimpl.Provide(istorage.ProvideMem())
	storage, _ := sp.AppStorage(istructs.AppQName_test1_app1)

	versions := vers.New()
	if err := versions.Prepare(storage); err != nil {
		panic(err)
	}

	testName := "test"
	appDefBuilder := appdef.New()
	appDefBuilder.AddStruct(appdef.NewQName("test", "el"), appdef.DefKind_Element).
		AddContainer(testName, appdef.NewQName("test", "el"), 0, appdef.Occurs_Unbounded)
	appDef, err := appDefBuilder.Build()
	if err != nil {
		panic(err)
	}

	containers := New()
	if err := containers.Prepare(storage, versions, appDef); err != nil {
		panic(err)
	}

	require := require.New(t)
	t.Run("basic Containers methods", func(t *testing.T) {
		id, err := containers.GetID(testName)
		require.NoError(err)
		require.NotEqual(NullContainerID, id)

		n, err := containers.GetContainer(id)
		require.NoError(err)
		require.Equal(testName, n)

		t.Run("must be able to load early stored names", func(t *testing.T) {
			otherVersions := vers.New()
			if err := otherVersions.Prepare(storage); err != nil {
				panic(err)
			}

			otherContainers := New()
			if err := otherContainers.Prepare(storage, versions, nil); err != nil {
				panic(err)
			}

			id1, err := containers.GetID(testName)
			require.NoError(err)
			require.Equal(id, id1)

			n1, err := containers.GetContainer(id)
			require.NoError(err)
			require.Equal(testName, n1)
		})
	})
}
