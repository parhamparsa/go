package persist

import (
	"app/internal/domain/entity"
	interfaces "app/internal/domain/interface"
	"app/pkg/sqlt"
	"context"
)

type transferRepository struct {
	store *MysqlStore
	table *sqlt.SqlTable
}

var _ interfaces.TransferRepository = transferRepository{}

func newTransferRepo(s *MysqlStore) interfaces.TransferRepository {
	//todo add created_at
	return &transferRepository{
		store: s,
		table: sqlt.NewTable("transfers", "ID", []string{
			"id", "from", "to", "amount"},
		),
	}
}

func (r transferRepository) AddTransferRecord(ctx context.Context, transfer *entity.Transfer) error {
	query, args, err := r.store.qb().
		Insert(r.table.Name()).
		Columns("`from`", "`to`", "amount").
		Values(transfer.From, transfer.To, transfer.Amount).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.masterDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
