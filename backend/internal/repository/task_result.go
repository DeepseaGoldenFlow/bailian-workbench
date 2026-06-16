package repository

func UpdateTaskResult(taskID, status, resultJSON string) {
	DB.Exec("UPDATE generation_history SET status = ?, result_urls = ?, updated_at = CURRENT_TIMESTAMP WHERE task_id = ?", status, resultJSON, taskID)
}
