package repository

import (
	"go.uber.org/mock/gomock"
	"testing"
	"ytb-video-sharing-app-be/mocks"
)

type testConfig struct {
	ctr *gomock.Controller
	db  *mocks.MockDatabase
	tx  *mocks.MockTx
	row *mocks.MockRow
}

func SetupTest(t *testing.T) *testConfig {
	ctr := gomock.NewController(t)
	db := mocks.NewMockDatabase(ctr)
	tx := mocks.NewMockTx(ctr)
	row := mocks.NewMockRow(ctr)

	return &testConfig{ctr: ctr, db: db, tx: tx, row: row}
}

func (cfg *testConfig) TearDownTest() {
	cfg.ctr.Finish()
}
