# RAG

RAG (Retrieval-Augmented Generation) agent for searching relevant documents from a knowledge base.

## Description

This agent indexes documents from `database.xlsx` on startup and uses vector search to find the most relevant documents based on user queries.

## Architecture

1. **Indexing** (on startup):
   - Reads documents from `database.xlsx`
   - Generates embeddings for each document
   - Stores them in an in-memory vector store

2. **Retrieval** (on request):
   - Generates an embedding for the user query
   - Finds Top-K most similar documents (cosine similarity)
   - Returns found documents with relevance scores

## Configuration

### Embedding Provider

The system uses **OpenRouter** for embedding generation:
- Model: `openai/text-embedding-3-small`
- Dimension: 1536
- API: `https://openrouter.ai/api/v1/embeddings`

### Settings

In `settings.json`:
- `top_k = 5` - number of documents to return
- `embedding_dimension = 1536` - embedding vector dimension
- `openrouter_embedding_url` - OpenRouter API endpoint
- `openrouter_embedding_model` - embedding model name

## Setup

### Database Preparation

The database must be in Excel format (`database.xlsx`). Documents should be in the first column (column A) of the first sheet.

#### Excel File Format

- First sheet (Sheet1)
- Column A contains documents (one per row)
- Empty rows are skipped


## Project Structure

```
rag-platform/
├── main.go                 # Entry point, indexing on startup
├── methods/
│   ├── send_message.go     # Request handling, document search
│   └── get_task.go         # Task status retrieval
├── rag/
│   ├── embedding.go        # Embedding generation (OpenRouter)
│   ├── vector_store.go     # In-memory vector store
│   ├── indexer.go          # Document indexing
│   └── retriever.go        # Relevant document retrieval
├── settings/
│   └── settings.go         # Settings parsing
├── settings.json           # Service configuration
└── database.xlsx           # Knowledge base (Excel format)
```

## Features

- **Vector Search**: Uses cosine similarity to find relevant documents
- **In-memory Storage**: Fast search without external dependencies
- **OpenRouter Embeddings**: Uses OpenRouter API for high-quality embedding generation
- **Automatic Indexing**: Documents are indexed on service startup
- **Excel Support**: Reads documents from Excel format
