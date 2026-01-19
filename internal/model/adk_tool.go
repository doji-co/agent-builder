package model

type AdkTool string

const (
	GOOGLE_SEARCH         AdkTool = "google_search"
	GOOGLE_CODE_EXECUTION AdkTool = "google_code_execution"

	VERTEX_AI_RAG_RETRIEVAL AdkTool = "vertex_ai_rag_retrieval"
	VERTEX_AI_SEARCH        AdkTool = "vertex_ai_search"

	BIGQUERY_LIST_DATASET_IDS   AdkTool = "bigquery_list_dataset_ids"
	BIGQUERY_GET_DATASET_INFO   AdkTool = "bigquery_get_dataset_info"
	BIGQUERY_LIST_TABLE_IDS     AdkTool = "bigquery_list_table_ids"
	BIGQUERY_GET_TABLE_INFO     AdkTool = "bigquery_get_table_info"
	BIGQUERY_EXECUTE_SQL        AdkTool = "bigquery_execute_sql"
	BIGQUERY_FORECAST           AdkTool = "bigquery_forecast"
	BIGQUERY_ASK_DATA_INSIGHTS  AdkTool = "bigquery_ask_data_insights"

	SPANNER_LIST_TABLE_NAMES   AdkTool = "spanner_list_table_names"
	SPANNER_GET_TABLE_SCHEMA   AdkTool = "spanner_get_table_schema"
	SPANNER_EXECUTE_SQL        AdkTool = "spanner_execute_sql"
	SPANNER_SIMILARITY_SEARCH  AdkTool = "spanner_similarity_search"

	BIGTABLE_LIST_INSTANCES    AdkTool = "bigtable_list_instances"
	BIGTABLE_GET_INSTANCE_INFO AdkTool = "bigtable_get_instance_info"
	BIGTABLE_LIST_TABLES       AdkTool = "bigtable_list_tables"
	BIGTABLE_GET_TABLE_INFO    AdkTool = "bigtable_get_table_info"
	BIGTABLE_EXECUTE_SQL       AdkTool = "bigtable_execute_sql"

	GKE_CODE_EXECUTOR AdkTool = "gke_code_executor"
)

func GetAllAdkTools() []AdkTool {
	return []AdkTool{
		GOOGLE_SEARCH,
		GOOGLE_CODE_EXECUTION,
		VERTEX_AI_RAG_RETRIEVAL,
		VERTEX_AI_SEARCH,
		BIGQUERY_LIST_DATASET_IDS,
		BIGQUERY_GET_DATASET_INFO,
		BIGQUERY_LIST_TABLE_IDS,
		BIGQUERY_GET_TABLE_INFO,
		BIGQUERY_EXECUTE_SQL,
		BIGQUERY_FORECAST,
		BIGQUERY_ASK_DATA_INSIGHTS,
		SPANNER_LIST_TABLE_NAMES,
		SPANNER_GET_TABLE_SCHEMA,
		SPANNER_EXECUTE_SQL,
		SPANNER_SIMILARITY_SEARCH,
		BIGTABLE_LIST_INSTANCES,
		BIGTABLE_GET_INSTANCE_INFO,
		BIGTABLE_LIST_TABLES,
		BIGTABLE_GET_TABLE_INFO,
		BIGTABLE_EXECUTE_SQL,
		GKE_CODE_EXECUTOR,
	}
}

func GetAdkToolDescription(tool AdkTool) string {
	descriptions := map[AdkTool]string{
		GOOGLE_SEARCH:               "Web searches using Google Search",
		GOOGLE_CODE_EXECUTION:       "Execute code for calculations and data manipulation",
		VERTEX_AI_RAG_RETRIEVAL:     "Private data retrieval using Vertex AI RAG Engine",
		VERTEX_AI_SEARCH:            "Search across private data stores via Vertex AI",
		BIGQUERY_LIST_DATASET_IDS:   "List BigQuery dataset IDs",
		BIGQUERY_GET_DATASET_INFO:   "Get BigQuery dataset information",
		BIGQUERY_LIST_TABLE_IDS:     "List BigQuery table IDs",
		BIGQUERY_GET_TABLE_INFO:     "Get BigQuery table information",
		BIGQUERY_EXECUTE_SQL:        "Execute SQL queries on BigQuery",
		BIGQUERY_FORECAST:           "BigQuery forecasting capabilities",
		BIGQUERY_ASK_DATA_INSIGHTS:  "Ask data insights from BigQuery",
		SPANNER_LIST_TABLE_NAMES:    "List Cloud Spanner table names",
		SPANNER_GET_TABLE_SCHEMA:    "Get Cloud Spanner table schema",
		SPANNER_EXECUTE_SQL:         "Execute SQL queries on Cloud Spanner",
		SPANNER_SIMILARITY_SEARCH:   "Similarity search on Cloud Spanner",
		BIGTABLE_LIST_INSTANCES:     "List Bigtable instances",
		BIGTABLE_GET_INSTANCE_INFO:  "Get Bigtable instance information",
		BIGTABLE_LIST_TABLES:        "List Bigtable tables",
		BIGTABLE_GET_TABLE_INFO:     "Get Bigtable table information",
		BIGTABLE_EXECUTE_SQL:        "Execute SQL queries on Bigtable",
		GKE_CODE_EXECUTOR:           "Secure code execution in GKE sandboxed environments",
	}
	return descriptions[tool]
}
