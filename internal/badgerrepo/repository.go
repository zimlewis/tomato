package badgerrepo

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/zimlewis/tomato/internal/tomatoerrs"
)

// Repository that connect to the Badger database
type Repository struct {
	db *badger.DB
}


var timerKey = []byte("timer")
var startTimeKey = []byte("start")

func New(db *badger.DB) Repository {
	return Repository{ 
		db: db,
	}
}

// Delete start time if it present
func (repo *Repository) DeleteStartTime(ctx context.Context) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(startTimeKey)
	})
	if err != nil {
		return errors.Join(tomatoerrs.ErrBadgerDB, err)
	}

	return nil
}

// Set start time to time passed in
func (repo *Repository) SetStartTime(ctx context.Context, time int64) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(startTimeKey, binary.BigEndian.AppendUint64(nil, uint64(time)))
		return err
	})
	if err != nil {
		return errors.Join(tomatoerrs.ErrBadgerDB, err)
	}

	return nil
}

// Update the clock by the clock index passed in
func (repo *Repository) SetClock(ctx context.Context, clockIndex int) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(timerKey, binary.BigEndian.AppendUint16(nil, uint16(clockIndex)))
		return err
	})
	if err != nil {
		return errors.Join(tomatoerrs.ErrBadgerDB, err)
	}

	return nil
}

// Will set the clock to pomodoro by default if it hasn't been set
func (repo *Repository) GetClock(ctx context.Context) (uint16, error) {
	var clock uint16
	err := repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(timerKey)
		if err == badger.ErrKeyNotFound {
			// Set clock to default value(0, pomodoro)
			err := repo.SetClock(ctx, 0)
			if err != nil { return err }
			clock = 0
			return nil
		} else if err != nil { return err }

		b, err := item.ValueCopy(nil)
		if err != nil { return err }

		clock = binary.BigEndian.Uint16(b)

		return nil
	})
	if err != nil {
		return 0, errors.Join(tomatoerrs.ErrBadgerDB, err)
	}

	return clock, nil
}

// Get the start time of the current session, if the session is not yet started, return a tomatoerrs.ErrDidNotStart
func (repo *Repository) GetStartTime(ctx context.Context) (int64, error) {
	var startTime int64

	err := repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(startTimeKey)
		if err != nil {
			return err
		}

		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		startTime = int64(binary.BigEndian.Uint64(b))

		return nil
	})
	
	if err == badger.ErrKeyNotFound {
		return 0, tomatoerrs.ErrDidNotStart
	}
	if err != nil {
		return 0, errors.Join(tomatoerrs.ErrBadgerDB, err)
	}


	return startTime, nil
}
