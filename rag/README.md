# RAG API Service

This directory contains the code for the Retrieval-Augmented Generation (RAG) API service.

## Running the Server

There are two modes to run the server:

### 1. Normal Mode (Database Based)

This mode requires a running Spanner database.

```bash
make run-history-api
```

This command uses `./configs/chrome-internal.json` by default.

### 2. In-Memory Mode

This mode loads the index snapshot directly from GCS and serves queries from memory.
It does not require a running Spanner database.

```bash
make run-history-api-in-memory
```

By default, it uses the date `2026-04-20`. You can override the date by setting the
`INDEX_DATE` environment variable:

```bash
INDEX_DATE=2026-04-20 make run-history-api-in-memory
```

**Note**: Ensure that your configuration file (e.g., `./configs/demo.json`) has the
`gcs_bucket` field set.

## Configuration

Configuration files are located in the `configs` directory. Key parameters include:

- `spanner_config`: For database mode.
- `gcs_bucket`: For in-memory mode.
- `query_embedding_model`: Gemini model for embeddings.
- `summary_model`: Gemini model for summaries.

## Environment Variables

Ensure you have the following environment variables set for Gemini API access:

- `GEMINI_API_KEY` (if running locally with an API key)
- Or `GEMINI_PROJECT` and `GEMINI_LOCATION` if using Google Cloud project credentials.
