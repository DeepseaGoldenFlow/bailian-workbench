package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.bailian.workbench.dto.image.ImageGenRequest;
import com.bailian.workbench.dto.image.ImageGenResponse;
import com.bailian.workbench.dto.task.TaskResponse;
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
public class ImageGenService {

    private static final Logger log = LoggerFactory.getLogger(ImageGenService.class);
    private static final int MAX_POLL_ATTEMPTS = 120;
    private static final long POLL_INTERVAL_SECONDS = 5;

    private final RestClient restClient;
    private final WebClient webClient;
    private final TaskService taskService;
    private final BailianConfig bailianConfig;

    public ImageGenService(RestClient restClient, WebClient webClient, TaskService taskService, BailianConfig bailianConfig) {
        this.restClient = restClient;
        this.webClient = webClient;
        this.taskService = taskService;
        this.bailianConfig = bailianConfig;
    }

    public String generateImage(ImageGenRequest request) {
        ImageGenResponse submitResponse = restClient.post()
                .uri("/api/v1/services/aigc/text2image/image-synthesis")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(ImageGenResponse.class);

        String taskId = submitResponse.output().taskId();
        log.info("Image generation task submitted: {}", taskId);

        // 2. Poll
        ImageGenResponse.Output result = pollTask(taskId);

        // 3. Download
        return downloadImages(result);
    }

    public ImageGenResponse submitTask(ImageGenRequest request) {
        return restClient.post()
                .uri("/api/v1/services/aigc/text2image/image-synthesis")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(ImageGenResponse.class);
    }

    private ImageGenResponse.Output pollTask(String taskId) {
        int attempts = 0;
        while (attempts < MAX_POLL_ATTEMPTS) {
            try {
                TaskResponse taskResponse = taskService.getTaskStatus(taskId);
                String status = taskResponse.output().taskStatus();

                if ("SUCCEEDED".equalsIgnoreCase(status)) {
                    return convertOutput(taskResponse.output());
                } else if ("FAILED".equalsIgnoreCase(status)) {
                    throw new RuntimeException("Task failed");
                }
                TimeUnit.SECONDS.sleep(POLL_INTERVAL_SECONDS);
                attempts++;
            } catch (Exception e) {
                throw new RuntimeException("Polling error: " + e.getMessage(), e);
            }
        }
        throw new RuntimeException("Timed out");
    }

    private ImageGenResponse.Output convertOutput(TaskResponse.Output output) {
        var results = output.results() != null
                ? output.results().stream().map(r -> new ImageGenResponse.Output.Result(r.url(), null, null)).toList()
                : java.util.List.<ImageGenResponse.Output.Result>of();
        return new ImageGenResponse.Output(output.taskId(), output.taskStatus(), results);
    }

    private String downloadImages(ImageGenResponse.Output output) {
        if (output.results() == null || output.results().isEmpty()) throw new RuntimeException("No results");

        String basePath = bailianConfig.getStorage().getBasePath();
        Path imagesDir = Path.of(basePath, "images");
        try { Files.createDirectories(imagesDir); } catch (IOException e) { throw new RuntimeException(e); }

        StringBuilder urls = new StringBuilder();
        for (var r : output.results()) {
            String filename = UUID.randomUUID() + ".png";
            Path localPath = imagesDir.resolve(filename);
            try {
                byte[] bytes = webClient.get().uri(r.url()).retrieve().onStatus(HttpStatusCode::isError, resp -> Mono.error(new RuntimeException("Download failed"))).bodyToMono(byte[].class).block();
                Files.write(localPath, bytes);
                if (urls.length() > 0) urls.append(",");
                urls.append("/storage/images/").append(filename);
            } catch (Exception e) { log.error("Download failed: {}", e.getMessage()); }
        }
        if (urls.isEmpty()) throw new RuntimeException("All downloads failed");
        return urls.toString();
    }
}
