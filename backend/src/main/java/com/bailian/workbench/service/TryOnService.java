package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.bailian.workbench.dto.task.TaskResponse;
import com.bailian.workbench.dto.tryon.TryOnRequest;
import com.bailian.workbench.dto.tryon.TryOnResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatusCode;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service
public class TryOnService {

    private static final Logger log = LoggerFactory.getLogger(TryOnService.class);
    private static final int MAX_POLL_ATTEMPTS = 120;
    private static final long POLL_INTERVAL_SECONDS = 5;

    private final RestClient restClient;
    private final WebClient webClient;
    private final TaskService taskService;
    private final BailianConfig bailianConfig;

    public TryOnService(RestClient restClient, WebClient webClient, TaskService taskService, BailianConfig bailianConfig) {
        this.restClient = restClient;
        this.webClient = webClient;
        this.taskService = taskService;
        this.bailianConfig = bailianConfig;
    }

    public String generateTryOn(TryOnRequest request) {
        TryOnResponse submitResponse = restClient.post()
                .uri("/api/v1/services/aigc/image2image/outfit-anyone")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(TryOnResponse.class);

        String taskId = submitResponse.output().taskId();
        log.info("AI Try-On task submitted: {}", taskId);

        TryOnResponse.Output result = pollTask(taskId);
        return downloadImage(result);
    }

    public TryOnResponse submitTask(TryOnRequest request) {
        return restClient.post()
                .uri("/api/v1/services/aigc/image2image/outfit-anyone")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(TryOnResponse.class);
    }

    private TryOnResponse.Output pollTask(String taskId) {
        int attempts = 0;
        while (attempts < MAX_POLL_ATTEMPTS) {
            try {
                TaskResponse taskResponse = taskService.getTaskStatus(taskId);
                String status = taskResponse.output().taskStatus();

                if ("SUCCEEDED".equalsIgnoreCase(status)) {
                    return convertOutput(taskResponse.output());
                } else if ("FAILED".equalsIgnoreCase(status)) {
                    throw new RuntimeException("Try-On task failed: " + taskResponse.message());
                }
                TimeUnit.SECONDS.sleep(POLL_INTERVAL_SECONDS);
                attempts++;
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new RuntimeException("Polling interrupted", e);
            } catch (Exception e) {
                throw new RuntimeException("Polling error: " + e.getMessage(), e);
            }
        }
        throw new RuntimeException("Try-On task timed out");
    }

    private TryOnResponse.Output convertOutput(TaskResponse.Output output) {
        String imageUrl = null;
        if (output.results() != null && !output.results().isEmpty()) {
            imageUrl = output.results().get(0).url();
        }
        return new TryOnResponse.Output(output.taskId(), output.taskStatus(), imageUrl);
    }

    private String downloadImage(TryOnResponse.Output output) {
        if (output.imageUrl() == null) throw new RuntimeException("No image result");

        String basePath = bailianConfig.getStorage().getBasePath();
        Path imagesDir = Path.of(basePath, "images");
        try { Files.createDirectories(imagesDir); } catch (IOException e) { throw new RuntimeException(e); }

        String filename = UUID.randomUUID() + ".png";
        Path localPath = imagesDir.resolve(filename);
        try {
            byte[] bytes = webClient.get().uri(output.imageUrl())
                    .retrieve()
                    .onStatus(HttpStatusCode::isError, resp -> Mono.error(new RuntimeException("Download failed")))
                    .bodyToMono(byte[].class).block();
            Files.write(localPath, bytes);
            return "/storage/images/" + filename;
        } catch (Exception e) {
            throw new RuntimeException("Image download failed: " + e.getMessage());
        }
    }
}
