package model

import "testing"

func TestAdkToolConstants(t *testing.T) {
	tests := []struct {
		name     string
		tool     AdkTool
		expected string
	}{
		{
			name:     "google search tool",
			tool:     GOOGLE_SEARCH,
			expected: "google_search",
		},
		{
			name:     "code execution tool",
			tool:     GOOGLE_CODE_EXECUTION,
			expected: "google_code_execution",
		},
		{
			name:     "vertex ai rag retrieval tool",
			tool:     VERTEX_AI_RAG_RETRIEVAL,
			expected: "vertex_ai_rag_retrieval",
		},
		{
			name:     "vertex ai search tool",
			tool:     VERTEX_AI_SEARCH,
			expected: "vertex_ai_search",
		},
		{
			name:     "bigquery list datasets tool",
			tool:     BIGQUERY_LIST_DATASET_IDS,
			expected: "bigquery_list_dataset_ids",
		},
		{
			name:     "bigquery get dataset info tool",
			tool:     BIGQUERY_GET_DATASET_INFO,
			expected: "bigquery_get_dataset_info",
		},
		{
			name:     "bigquery list tables tool",
			tool:     BIGQUERY_LIST_TABLE_IDS,
			expected: "bigquery_list_table_ids",
		},
		{
			name:     "bigquery get table info tool",
			tool:     BIGQUERY_GET_TABLE_INFO,
			expected: "bigquery_get_table_info",
		},
		{
			name:     "bigquery execute sql tool",
			tool:     BIGQUERY_EXECUTE_SQL,
			expected: "bigquery_execute_sql",
		},
		{
			name:     "bigquery forecast tool",
			tool:     BIGQUERY_FORECAST,
			expected: "bigquery_forecast",
		},
		{
			name:     "bigquery ask data insights tool",
			tool:     BIGQUERY_ASK_DATA_INSIGHTS,
			expected: "bigquery_ask_data_insights",
		},
		{
			name:     "spanner list tables tool",
			tool:     SPANNER_LIST_TABLE_NAMES,
			expected: "spanner_list_table_names",
		},
		{
			name:     "spanner get table schema tool",
			tool:     SPANNER_GET_TABLE_SCHEMA,
			expected: "spanner_get_table_schema",
		},
		{
			name:     "spanner execute sql tool",
			tool:     SPANNER_EXECUTE_SQL,
			expected: "spanner_execute_sql",
		},
		{
			name:     "spanner similarity search tool",
			tool:     SPANNER_SIMILARITY_SEARCH,
			expected: "spanner_similarity_search",
		},
		{
			name:     "bigtable list instances tool",
			tool:     BIGTABLE_LIST_INSTANCES,
			expected: "bigtable_list_instances",
		},
		{
			name:     "bigtable get instance info tool",
			tool:     BIGTABLE_GET_INSTANCE_INFO,
			expected: "bigtable_get_instance_info",
		},
		{
			name:     "bigtable list tables tool",
			tool:     BIGTABLE_LIST_TABLES,
			expected: "bigtable_list_tables",
		},
		{
			name:     "bigtable get table info tool",
			tool:     BIGTABLE_GET_TABLE_INFO,
			expected: "bigtable_get_table_info",
		},
		{
			name:     "bigtable execute sql tool",
			tool:     BIGTABLE_EXECUTE_SQL,
			expected: "bigtable_execute_sql",
		},
		{
			name:     "gke code executor tool",
			tool:     GKE_CODE_EXECUTOR,
			expected: "gke_code_executor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.tool) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.tool))
			}
		})
	}
}

func TestAdkToolString(t *testing.T) {
	tool := GOOGLE_SEARCH
	expected := "google_search"

	if string(tool) != expected {
		t.Errorf("expected %s, got %s", expected, string(tool))
	}
}
