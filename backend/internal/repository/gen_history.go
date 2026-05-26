package repository

import (
	"database/sql"
	"time"
)

func SaveGeneration(genType, model, prompt, resultJSON, taskID, status string) {
	DB.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, task_id, status) VALUES (?, ?, ?, ?, ?, ?)",
		genType, model, prompt, resultJSON, taskID, status)
}

func GetGenerationHistory(genType string) ([]map[string]any, error) {
	query := "SELECT id, type, model, prompt, result_urls, task_id, status, created_at FROM generation_history"
	args := []any{}
	if genType != "" {
		query += " WHERE type = ?"
		args = append(args, genType)
	}
	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]map[string]any, 0)
	for rows.Next() {
		var id int64
		var typ, model, prompt, resultJSON, taskID, status string
		var createdAt time.Time
		var promptNull, taskNull, resultNull sql.NullString
		rows.Scan(&id, &typ, &model, &promptNull, &resultNull, &taskNull, &status, &createdAt)
		prompt = promptNull.String
		resultJSON = resultNull.String
		taskID = taskNull.String
		history = append(history, map[string]any{
			"id":         id,
			"type":       typ,
			"model":      model,
			"prompt":     prompt,
			"result":     resultJSON,
			"task_id":    taskID,
			"status":     status,
			"created_at": createdAt,
		})
	}
	return history, nil
}

func DeleteGeneration(id string) {
	DB.Exec("DELETE FROM generation_history WHERE id = ?", id)
}