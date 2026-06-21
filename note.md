> Postgres error code

| Error Code | Error Name | Common Cause |
| --- | --- | --- |
| 23505 | UNIQUE_VIOLATION | Trying to insert a duplicate value into a primary key or unique index column. |
| 23503 | FOREIGN_KEY_VIOLATION | Inserting a reference ID that does not exist in the parent table. |
| 23502 | NOT_NULL_VIOLATION | Leaving a mandatory field blank (trying to insert a NULL value). |
| 42601 | SYNTAX_ERROR | Misplaced commas, misspelled keywords, or unmatched parentheses in your SQL. |
| 42703 | UNDEFINED_COLUMN | Typing a column name that does not exist in the target table. |
| 42P01 | UNDEFINED_TABLE | Querying a table name that is misspelled or does not exist. |
| 28P01 | INVALID_PASSWORD | Providing the incorrect password during application database authentication. |
| 40001 | SERIALIZATION_FAILURE | Concurrent transactions conflicting under the serializable isolation level. |
| 57014 | QUERY_CANCELED | The statement was intentionally stopped or hit a configured statement timeout. |

> Major Issue

| Class Code | Category Name | Description |
| --- | --- | --- |
| 08 | Connection Exception | Database connectivity failures or timeout drops. |
| 22 | Data Exception | Math anomalies (e.g., division by zero), string truncation, or bad date formats. |
| 25 | Invalid Transaction State | Trying to run commands inside a broken transaction block. |
| 42 | Syntax Error or Access Rule Violation | Query structural errors or missing table permissions. |
| 53 | Insufficient Resources | Disk space exhaustion or memory allocation issues. |
