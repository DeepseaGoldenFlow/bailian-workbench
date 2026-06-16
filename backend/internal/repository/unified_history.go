package repository

import (
	"database/sql"
	"time"
)

type UnifiedEntry struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Content   string `json:"content"`
	Result    string `json:"result"`
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func GetUnifiedHistory(typ string, limit int) ([]UnifiedEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	query := `
		SELECT id, 'chat' AS type, '' AS prompt, model, content, '' AS result, '' AS task_id, '' AS status, created_at
		FROM chat_history WHERE role = 'user'
		UNION ALL
		SELECT id, type, COALESCE(prompt,''), model, '', COALESCE(result_urls,''), COALESCE(task_id,''), COALESCE(status,''), created_at
		FROM generation_history
	`
	args := []any{}
	if typ != "" {
		if typ == "chat" {
			query = `
				SELECT id, 'chat' AS type, '' AS prompt, model, content, '' AS result, '' AS task_id, '' AS status, created_at
				FROM chat_history WHERE role = 'user'
			`
		} else {
			query = `
				SELECT id, type, COALESCE(prompt,''), model, '', COALESCE(result_urls,''), COALESCE(task_id,''), COALESCE(status,''), created_at
				FROM generation_history WHERE type = ?
			`
			args = append(args, typ)
		}
	}
	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]UnifiedEntry, 0)
	for rows.Next() {
		var e UnifiedEntry
		var createdAt time.Time
		var promptNull, contentNull, resultNull, taskNull, statusNull sql.NullString
		rows.Scan(&e.ID, &e.Type, &promptNull, &e.Model, &contentNull, &resultNull, &taskNull, &statusNull, &createdAt)
		e.Prompt = promptNull.String
		e.Content = contentNull.String
		e.Result = resultNull.String
		e.TaskID = taskNull.String
		e.Status = statusNull.String
		e.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		entries = append(entries, e)
	}
	return entries, nil
}

func DeleteUnifiedEntry(typ string, id string) error {
	if typ == "chat" {
		DB.Exec("DELETE FROM chat_history WHERE id = ?", id)
	} else {
		DB.Exec("DELETE FROM generation_history WHERE id = ?", id)
	}
	return nil
}
