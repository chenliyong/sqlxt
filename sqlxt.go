package sqlxt

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Sqlxt struct {
	db    *sqlx.DB
	tx    *sqlx.Tx
	query *Query
}

func New(db *sqlx.DB, q *Query) *Sqlxt {
	return &Sqlxt{
		db:    db,
		query: q,
	}
}

func NewTx(tx *sqlx.Tx, q *Query) *Sqlxt {
	return &Sqlxt{
		tx:    tx,
		query: q,
	}
}

func (st *Sqlxt) Get(dest interface{}) error {
	st.query.Limit(1)
	query, args, err := st.query.BuildQuery()
	if err != nil {
		return err
	}
	var row *sqlx.Row
	if st.tx != nil {
		row = st.tx.Unsafe().QueryRowx(query, args...)
	} else {
		row = st.db.Unsafe().QueryRowx(query, args...)
	}
	return row.StructScan(dest)
}

func (st *Sqlxt) All(dest interface{}) error {
	var err error
	query, args, err := st.query.BuildQuery()
	if err != nil {
		return err
	}
	var rows *sqlx.Rows
	if st.tx != nil {
		rows, err = st.tx.Unsafe().Queryx(query, args...)
	} else {
		rows, err = st.db.Unsafe().Queryx(query, args...)
	}
	if err != nil {
		return err
	}
	return rows.StructScan(dest)
}

func (st *Sqlxt) Update(data map[string]interface{}) (sql.Result, error) {
	return st.Exec("update", data)
}

func (st *Sqlxt) Insert(data map[string]interface{}) (sql.Result, error) {
	return st.Exec("insert", data)
}

func (st *Sqlxt) Delete() (sql.Result, error) {
	return st.Exec("delete", nil)
}

func (st *Sqlxt) Exec(method string, data map[string]interface{}) (sql.Result, error) {
	var err error
	query, args, err := st.query.BuildExec(method, data)
	if err != nil {
		return nil, err
	}
	var result sql.Result
	if st.tx != nil {
		result, err = st.tx.Exec(query, args...)
	} else {
		result, err = st.db.Exec(query, args...)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
