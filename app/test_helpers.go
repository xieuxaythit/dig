package app

import (
	"encoding/json"
	"os"

	"github.com/CosmWasm/wasmd/app"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appv1 "github.com/notional-labs/dig/app"
	"github.com/tendermint/spm/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// Setup initializes a new DigApp
func Setup(isCheckTx bool) cosmoscmd.App {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)

	digapp := NewDigApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(encCdc.Marshaler)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		digapp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return digapp
}

func SetupAppv1(isCheckTx bool) cosmoscmd.App {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)

	digapp := appv1.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})

	if !isCheckTx {
		genesisState := NewDefaultGenesisState(encCdc.Marshaler)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		digapp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return digapp
}

// SetupTestingAppWithLevelDb initializes a new App intended for testing,
// with LevelDB as a db
func SetupTestingAppWithLevelDb(isCheckTx bool, dir string) (digapp cosmoscmd.App, cleanupFn func()) {
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)

	db, err := sdk.NewLevelDB("dig_leveldb_testing", dir)
	if err != nil {
		panic(err)
	}
	digapp = NewDigApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(encCdc.Marshaler)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		digapp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	cleanupFn = func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	return
}
