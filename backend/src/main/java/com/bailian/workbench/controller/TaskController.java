package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.task.TaskResponse;
import com.bailian.workbench.service.TaskService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/tasks")
public class TaskController {

    private static final Logger log = LoggerFactory.getLogger(TaskController.class);

    private final TaskService taskService;

    public TaskController(TaskService taskService) {
        this.taskService = taskService;
    }

    /**
     * Get the status of a task by its ID.
     */
    @GetMapping("/{task_id}")
    public ResponseEntity<?> getTaskStatus(@PathVariable("task_id") String taskId) {
        try {
            log.info("Fetching task status for: {}", taskId);
            TaskResponse response = taskService.getTaskStatus(taskId);
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            log.error("Failed to get task status for {}: {}", taskId, e.getMessage(), e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(
                            HttpStatus.INTERNAL_SERVER_ERROR.value(),
                            "Failed to get task status: " + e.getMessage(),
                            null
                    ));
        }
    }
}
