package index_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/mitchellh/hashstructure"
	"github.com/stretchr/testify/require"

	"github.com/ovh/cds/engine/cdn/index"
	cdntest "github.com/ovh/cds/engine/cdn/test"
	"github.com/ovh/cds/engine/gorpmapper"
	"github.com/ovh/cds/engine/test"
	"github.com/ovh/cds/sdk"
)

func TestLoadItem(t *testing.T) {
	m := gorpmapper.New()
	index.InitDBMapping(m)

	db, _ := test.SetupPGWithMapper(t, m, sdk.TypeCDN)
	cdntest.ClearIndex(t, context.TODO(), m, db)

	apiRef := index.ApiRef{
		ProjectKey: sdk.RandomString(10),
	}
	hashRefU, err := hashstructure.Hash(apiRef, nil)
	require.NoError(t, err)
	hashRef := strconv.FormatUint(hashRefU, 10)

	i := index.Item{
		ApiRef:     apiRef,
		ApiRefHash: hashRef,
		Type:       index.TypeItemStepLog,
	}
	require.NoError(t, index.InsertItem(context.TODO(), m, db, &i))
	t.Cleanup(func() { _ = index.DeleteItem(m, db, &i) })

	res, err := index.LoadItemByID(context.TODO(), m, db, i.ID)
	require.NoError(t, err)
	require.Equal(t, i.ID, res.ID)
	require.Equal(t, i.Type, res.Type)
}
