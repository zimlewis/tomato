package repository

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/dgraph-io/badger/v4"
	errs "github.com/zimlewis/tomato/internal/errors"
)

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

func (repo *Repository) DeleteStartTime(ctx context.Context) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(startTimeKey)
	})
	if err != nil {
		return errors.Join(errs.ErrBadgerDB, err)
	}

	return nil
}

func (repo *Repository) SetStartTime(ctx context.Context, time int64) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(startTimeKey, binary.BigEndian.AppendUint64(nil, uint64(time)))
		return err
	})
	if err != nil {
		return errors.Join(errs.ErrBadgerDB, err)
	}

	return nil
}

func (repo *Repository) SetClock(ctx context.Context, clockIndex int) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(timerKey, binary.BigEndian.AppendUint16(nil, uint16(clockIndex)))
		return err
	})
	if err != nil {
		return errors.Join(errs.ErrBadgerDB, err)
	}

	return nil
}

// Will set it to pomodoro by default if the clock haven't been set
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
		return 0, errors.Join(errs.ErrBadgerDB, err)
	}

	return clock, nil
}

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
		return 0, errs.ErrDidNotStart
	}
	if err != nil {
		return 0, errors.Join(errs.ErrBadgerDB, err)
	}


	return startTime, nil
}
