package blockchain

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Block struct {
	Index        int         `json:"index"`
	Timestamp    time.Time   `json:"timestamp"`
	Transaction  Transaction `json:"transaction"`
	PreviousHash string      `json:"previous_hash"`
	Hash         string      `json:"hash"`
}

type Transaction struct {
	Action    string    `json:"action"`     // "Create" or "Update"
	RecordID  int       `json:"record_id"`  // ID of the medical record
	DoctorID  int       `json:"doctor_id"`  // Doctor who performed the action
	PatientID int       `json:"patient_id"` // Patient associated with the record
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"` // JSON string of record data
}

type Blockchain struct {
	db    *sql.DB
	Chain []Block
}

func NewBlockchain(db *sql.DB) *Blockchain {
	bc := &Blockchain{db: db, Chain: make([]Block, 0)}
	bc.loadFromDB() // Load existing blocks from DB on initialization
	if len(bc.Chain) == 0 {
		bc.Chain = append(bc.Chain, GenesisBlock())
		bc.saveBlock(GenesisBlock()) // Save genesis block to DB
	}
	return bc
}

func GenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transaction:  Transaction{},
		PreviousHash: "0",
		Hash:         calculateHash(0, time.Now(), Transaction{}, "0"),
	}
}

func calculateHash(index int, timestamp time.Time, transaction Transaction, previousHash string) string {
	record := struct {
		Index        int
		Timestamp    time.Time
		Transaction  Transaction
		PreviousHash string
	}{
		Index:        index,
		Timestamp:    timestamp,
		Transaction:  transaction,
		PreviousHash: previousHash,
	}

	data, _ := json.Marshal(record)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func (bc *Blockchain) AddBlock(action string, recordID, doctorID, patientID int, details string) {
	previousBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:     previousBlock.Index + 1,
		Timestamp: time.Time{},
		Transaction: Transaction{
			Action:    action,
			RecordID:  recordID,
			DoctorID:  doctorID,
			PatientID: patientID,
			Timestamp: time.Now(),
			Details:   details,
		},
		PreviousHash: previousBlock.Hash,
		Hash:         "",
	}

	newBlock.Timestamp = time.Now()
	newBlock.Hash = calculateHash(newBlock.Index, newBlock.Timestamp, newBlock.Transaction, newBlock.PreviousHash)

	bc.Chain = append(bc.Chain, newBlock)
	bc.saveBlock(newBlock) // Save to DB
}

func (bc *Blockchain) saveBlock(block Block) error {
	transactionJSON, err := json.Marshal(block.Transaction)
	if err != nil {
		return err
	}

	_, err = bc.db.Exec(`
        INSERT INTO public.blocks (index, timestamp, transaction, previous_hash, hash)
        VALUES ($1, $2, $3, $4, $5)`,
		block.Index, block.Timestamp, transactionJSON, block.PreviousHash, block.Hash)
	return err
}

func (bc *Blockchain) loadFromDB() {
	rows, err := bc.db.Query("SELECT index, timestamp, transaction, previous_hash, hash FROM public.blocks ORDER BY index ASC")
	if err != nil {
		return // Handle error appropriately in production
	}
	defer rows.Close()

	for rows.Next() {
		var block Block
		var transactionJSON []byte
		err := rows.Scan(&block.Index, &block.Timestamp, &transactionJSON, &block.PreviousHash, &block.Hash)
		if err != nil {
			continue // Skip invalid blocks
		}

		var transaction Transaction
		if err := json.Unmarshal(transactionJSON, &transaction); err != nil {
			continue // Skip if transaction JSON is invalid
		}

		block.Transaction = transaction
		bc.Chain = append(bc.Chain, block)
	}
}
