package repository

import "time"

func SaveChatMessage(sessionID, role, model, content string) {
	DB.Exec("INSERT INTO chat_history (session_id, role, model, content) VALUES (?, ?, ?, ?)",
		sessionID, role, model, content)
}

func GetChatHistory() ([]map[string]any, error) {
	rows, err := DB.Query("SELECT id, session_id, role, model, content, created_at FROM chat_history ORDER BY created_at DESC LIMIT 200")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]map[string]any, 0)
	for rows.Next() {
		var id int64
		var sessionID, role, model, content string
		var createdAt time.Time
		rows.Scan(&id, &sessionID, &role, &model, &content, &createdAt)
		history = append(history, map[string]any{
			"id":         id,
			"session_id": sessionID,
			"role":       role,
			"model":      model,
			"content":    content,
			"created_at": createdAt,
		})
	}
	return history, nil
}

func DeleteChatMessage(id string) {
	DB.Exec("DELETE FROM chat_history WHERE id = ?", id)
}