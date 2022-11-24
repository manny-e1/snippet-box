package models

import "database/sql"

type TransactionExample struct {
	DB *sql.DB
}

func (trx *TransactionExample) InsertAndUpdate() error {
	tx, err := trx.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}

	result, err := tx.Exec(`INSERT INTO snippets(title,content,created,expires) 
VALUES("Transaction trail 2","let's if the transaction's working",UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 4 DAY))`)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	updStmt := `UPDATE snippets SET title="Transaction Trial 2" WHERE id=?`
	_, err = tx.Exec(updStmt, int(id))

	if err != nil {
		return err
	}
	err = tx.Commit()
	return err

}
