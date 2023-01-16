package utils

import (
	"context"
	"github.com/genjidb/genji"
)

func (u *Utils) ConnectGenji(dbFilePath string) (*genji.DB, error) {
	return genji.Open(dbFilePath)
}

func (u *Utils) ConnectGenjiWithContext(ctx context.Context, dbFilePath string) (*genji.DB, error) {
	db, dbErr := genji.Open(dbFilePath)
	if dbErr != nil {
		return nil, dbErr
	}

	return db.WithContext(ctx), nil
}

