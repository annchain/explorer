package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createBlockCollectSQL = `CREATE TABLE IF NOT EXISTS block_t
    (
		height				INTEGER NOT NULL,
		hash				VARCHAR(64) PRIMARY KEY,
		chain_id			VARCHAR(64) NOT NULL,
		time				DATETIME NOT NULL,
		num_txs				INTEGER,
		last_commit_hash    VARCHAR(64) NOT NULL,
		data_hash			VARCHAR(64) NOT NULL,
		validators_hash		VARCHAR(64) NOT NULL,
		app_hash			VARCHAR(64),
		reward				INTEGER,
		coin_base			VARCHAR(64)
	);`

	createBlockIndex = `CREATE INDEX  IF NOT EXISTS heightindex ON BLOCK_T(height) ;`

	createTxCollectSQL = `CREATE TABLE IF NOT EXISTS transaction_t
    (
		hash			VARCHAR(64),
		payload			VARCHAR(64) NOT NULL,
		payload_hex		VARCHAR(64) NOT NULL,
		from_addr		VARCHAR(64) NOT NULL,
		to_addr			VARCHAR(64) NOT NULL,
		receipt			VARCHAR(64) NOT NULL,
		amount			VARCHAR(64) NOT NULL,
		nonce			INTEGER NOT NULL,
		gas				VARCHAR(64) NOT NULL,
		size			INTEGER NOT NULL,
		block			VARCHAR(64) NOT NULL,
		contract		VARCHAR(64),
		time			DATETIME NOT NULL,
		height			INTEGER NOT NULL,
		tx_type			VARCHAR(64),
		fee				INTEGER
    );`

	createTrxIndex = `CREATE INDEX  IF NOT EXISTS hasindex ON transaction_t(hash);`


	blockSortSQL   = `SELECT * FROM block_t ORDER BY height DESC limit ?`
	blockRangeSQL  = `SELECT * FROM block_t WHERE height >= ? AND height <= ?`
	pageBlockSql   = `SELECT * FROM block_t ORDER BY height DESC limit ?,?`
	blockHashSQL   = `SELECT * FROM block_t WHERE hash = ?`
	blockHeightSQL = `SELECT * FROM block_t WHERE height = ?`
	blockInsertSQL = `INSERT INTO block_t VALUES(?,?,?,?,?,?,?,?,?,?,?)`

	txsByBlockHashSQL = `SELECT * FROM transaction_t WHERE block = ?`
	txsByToSQL        = `SELECT * FROM transaction_t WHERE to_addr = ?`
	txsByFromOrToSQL  = `SELECT * FROM transaction_t WHERE to_addr = ? OR from_addr = ?`
	txsNoContractSQL  = `SELECT * FROM transaction_t WHERE contract = "" ORDER BY height DESC limit ?`
	txsContractSQL    = `SELECT * FROM transaction_t WHERE contract != "" limit ?`
	txHashSQL         = `SELECT * FROM transaction_t WHERE hash = ?`
	txContractSQL     = `SELECT * FROM transaction_t WHERE contract = ?`
	txInsertSQL       = `INSERT INTO transaction_t VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	countSQL      = "SELECT COUNT(hash) FROM %s"
	dbBusyTimeout = "30000"
)

var db *sql.DB

func CreateSqlite() (err error) {
	db, err = sql.Open("sqlite3", fmt.Sprintf("./%s.db?cache=shared&_busy_timeout=%s", DB_NAME, dbBusyTimeout))
	if err != nil {
		return
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(createBlockCollectSQL)
	if err != nil {
		return
	}
	_, err = db.Exec(createTxCollectSQL)
	_, err = db.Exec(createBlockIndex)
	_, err = db.Exec(createTrxIndex)

	return
}

func LatestBlocksBySqlite(limit int) (blocks []Block, err error) {
	rows, errQ := db.Query(blockSortSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var block Block
		err = rows.Scan(&block.Height, &block.Hash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.Reward, &block.CoinBase)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}
	err = rows.Err()
	return
}

func SaveBlockBySqlite(blocks []Block) (err error) {

	sqlTx, errB := db.Begin()
	if errB != nil {
		err = errB
		return
	}
	stmt, errP := sqlTx.Prepare(blockInsertSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	for _, block := range blocks {
		_, err = stmt.Exec(
			block.Height,
			block.Hash,
			block.ChainID,
			block.Time,
			block.NumTxs,
			block.LastCommitHash,
			block.DataHash,
			block.ValidatorsHash,
			block.AppHash,
			block.Reward,
			block.CoinBase,
		)
		if err != nil {
			return
		}
	}
	err = sqlTx.Commit()
	return
}

func SaveTxBySqlite(txs []Transaction) (err error) {

	sqlTx, errB := db.Begin()
	if errB != nil {
		err = errB
		return
	}
	stmt, errP := sqlTx.Prepare(txInsertSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	for _, tx := range txs {
		_, err = stmt.Exec(
			tx.Hash,
			tx.Payload,
			tx.PayloadHex,
			tx.From,
			tx.To,
			tx.Receipt,
			tx.Amount,
			tx.Nonce,
			tx.Gas,
			tx.Size,
			tx.Block,
			tx.Contract,
			tx.Time,
			tx.Height,
			tx.TxType,
			tx.Fee,
		)
		if err != nil {
			return
		}
	}
	err = sqlTx.Commit()
	return
}

func PageBlocksBySqlite(pageIndex, pageSize int) (blocks []Block, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	rows, errQ := db.Query(pageBlockSql, (pageIndex-1)*pageSize, pageSize)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var block Block
		err = rows.Scan(&block.Height, &block.Hash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.Reward, &block.CoinBase)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}
	err = rows.Err()
	return
}

func BlocksFromToBySqlite(from, to int) (blocks []Block, err error) {

	rows, errQ := db.Query(blockRangeSQL, from, to)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var block Block
		err = rows.Scan(&block.Height, &block.Hash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.Reward, &block.CoinBase)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}
	err = rows.Err()
	return
}
func OneBlockBySqlite(hash string) (block Block, txs []Transaction, err error) {

	stmt, errP := db.Prepare(blockHashSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&block.Height, &block.Hash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.Reward, &block.CoinBase)
	if err != nil {
		return
	}
	txs, err = TransactionsByBlkhashBySqlite(hash)
	return
}

func BlockByHeightBySqlite(height int) (block Block, err error) {

	stmt, errP := db.Prepare(blockHeightSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(height).Scan(&block.Height, &block.Hash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.Reward, &block.CoinBase)
	return
}

func TransactionFromToBySqlite(from, to int) (txs []Transaction, err error) {
	return
}

func TransactionsByBlkhashBySqlite(hash string) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsByBlockHashSQL, hash)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func OneTransactionBySqlite(hash string) (tx Transaction, err error) {

	stmt, errP := db.Prepare(txHashSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
	return
}

func OneContractBySqlite(hash string) (tx Transaction, txs []Transaction, err error) {

	stmt, errP := db.Prepare(txContractSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
	if err != nil {
		return
	}

	txs, err = txsByTo(hash)
	return
}

func TxsBySqlite(limit int) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsNoContractSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func TxsQueryBySqlite(fromTo string) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsByFromOrToSQL, fromTo)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func ContractsBySqlite(limit int) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsContractSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func ContractBySqlite(hash string) (tx Transaction, txs []Transaction, err error) {

	tx, err = OneTransactionBySqlite(hash)
	if err != nil {
		return
	}
	txs, err = txsByTo(hash)
	return
}
func txsByTo(hash string) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsByToSQL, hash)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height, &tx.TxType, &tx.Fee)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}
func CollectionItemNumBySqlite(collect string) (count int, err error) {

	stmt, errP := db.Prepare(fmt.Sprintf(countSQL, collect))
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&count)
	return
}
