package com.bailian.workbench.service;

import com.bailian.workbench.dto.task.TaskResponse;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

@Service
public class TaskService {

    private final RestClient restClient;

    public TaskService(RestClient restClient) {
        this.restClient = restClient;
    }

    public TaskResponse getTaskStatus(String taskId) {
        // Manually concatenate URI to ensure ID is passed correctly
        return restClient.get()
                .uri("/api/v1/tasks/" + taskId)
                .retrieve()
                .body(TaskResponse.class);
    }
}
