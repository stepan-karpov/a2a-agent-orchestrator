package rag

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// Indexer indexes documents from database.xlsx
type Indexer struct {
	embeddingClient *EmbeddingClient
	vectorStore     *VectorStore
}

// NewIndexer creates a new indexer
func NewIndexer(embeddingClient *EmbeddingClient, vectorStore *VectorStore) *Indexer {
	return &Indexer{
		embeddingClient: embeddingClient,
		vectorStore:     vectorStore,
	}
}

// IndexDatabase indexes documents from database.xlsx
func (idx *Indexer) IndexDatabase(dbPath string) error {
	// Открываем Excel файл
	f, err := excelize.OpenFile(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database.xlsx: %w", err)
	}
	defer f.Close()

	// Получаем имя первого листа
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("no sheets found in database.xlsx")
	}

	// Читаем все строки из первого столбца (колонка A)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows from database.xlsx: %w", err)
	}

	var documents []string
	for i, row := range rows {
		// Пропускаем пустые строки
		if len(row) == 0 || row[0] == "" {
			continue
		}
		// Берем первую колонку (A) как документ
		documents = append(documents, row[0])
		if (i+1)%100 == 0 {
			fmt.Printf("Loaded %d documents from Excel...\n", i+1)
		}
	}

	fmt.Printf("Indexing %d documents...\n", len(documents))

	// Index each document
	for i, doc := range documents {
		if doc == "" {
			continue
		}

		// Generate an embedding
		embedding, err := idx.embeddingClient.GetEmbedding(doc)
		if err != nil {
			return fmt.Errorf("failed to generate embedding for document %d: %w", i, err)
		}

		// Add to vector store
		idx.vectorStore.Add(doc, embedding)

		if (i+1)%10 == 0 {
			fmt.Printf("Indexed %d/%d documents\n", i+1, len(documents))
		}
	}

	fmt.Printf("Successfully indexed %d documents\n", len(documents))
	return nil
}

// GetDatabasePath returns the path to database.xlsx
func GetDatabasePath() string {
	possiblePaths := []string{
		"database.xlsx",
		"./database.xlsx",
		"services/rag/database.xlsx",
		"services/agents/rag/database.xlsx",
		"/home/ubuntu/a2a-agent-orchestrator/services/rag/database.xlsx",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	return "database.xlsx"
}
