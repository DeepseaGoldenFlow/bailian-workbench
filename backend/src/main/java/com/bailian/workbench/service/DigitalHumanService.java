package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.bailian.workbench.dto.task.TaskResponse;
import com.bailian.workbench.dto.video.VideoGenRequest;
import com.bailian.workbench.dto.video.VideoGenResponse;
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
public class DigitalHumanService {
    private static final Logger log = LoggerFactory.getLogger(DigitalHumanService.class);
    private static final int MAX_POLL_ATTEMPTS = 600;
    private static final long POLL_INTERVAL_SECONDS = 10;

    private final RestClient restClient;
    private final WebClient webClient;
    private final TaskService taskService;
    private final BailianConfig bailianConfig;

    public DigitalHumanService(RestClient restClient, WebClient webClient, TaskService taskService, BailianConfig bailianConfig) {
        this.restClient = restClient;
        this.webClient = webClient;
        this.taskService = taskService;
        this.bailianConfig = bailianConfig;
    }

    public String generateVideo(VideoGenRequest request) {
        // 1. Detect (Optional but recommended, simplified here to just generation)
        
        // 2. Generate
        VideoGenResponse submitResponse = restClient.post()
                .uri("/api/v1/services/aigc/image2video/s2v/generate") // Assuming standard S2V endpoint structure
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(VideoGenResponse.class);

        String taskId = submitResponse.output().taskId();
        log.info("Digital Human S2V task submitted: {}", taskId);
        return downloadVideo(pollTask(taskId));
    }

    public VideoGenResponse submitTask(VideoGenRequest request) {
        return restClient.post()
                .uri("/api/v1/services/aigc/image2video/s2v/generate")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(VideoGenResponse.class);
    }

    private VideoGenResponse.Output pollTask(String taskId) {
        int attempts = 0;
        while (attempts < MAX_POLL_ATTEMPTS) {
            try {
                TaskResponse taskResponse = taskService.getTaskStatus(taskId);
                String status = taskResponse.output().taskStatus();
                if ("SUCCEEDED".equalsIgnoreCase(status)) {
                    return convertOutput(taskResponse.output());
                } else if ("FAILED".equalsIgnoreCase(status)) {
                    throw new RuntimeException("Task failed: " + taskResponse.message());
                }
                TimeUnit.SECONDS.sleep(POLL_INTERVAL_SECONDS);
                attempts++;
            } catch (Exception e) {
                throw new RuntimeException("Polling error: " + e.getMessage(), e);
            }
        }
        throw new RuntimeException("Timed out");
    }

    private VideoGenResponse.Output convertOutput(TaskResponse.Output output) {
        String videoUrl = null;
        if (output.results() != null && !output.results().isEmpty()) {
            videoUrl = output.results().get(0).url();
        }
        return new VideoGenResponse.Output(output.taskId(), output.taskStatus(), videoUrl);
    }

    private String downloadVideo(VideoGenResponse.Output output) {
        if (output.videoUrl() == null) throw new RuntimeException("No video result");
        String basePath = bailianConfig.getStorage().getBasePath();
        Path videosDir = Path.of(basePath, "videos");
        try { Files.createDirectories(videosDir); } catch (IOException e) { throw new RuntimeException(e); }
        String filename = UUID.randomUUID() + ".mp4";
        Path localPath = videosDir.resolve(filename);
        try {
            byte[] bytes = webClient.get().uri(output.videoUrl()).retrieve()
                    .onStatus(HttpStatusCode::isError, resp -> Mono.error(new RuntimeException("Download failed")))
                    .bodyToMono(byte[].class).block();
            Files.write(localPath, bytes);
            return "/storage/videos/" + filename;
        } catch (Exception e) {
            throw new RuntimeException("Download failed: " + e.getMessage());
        }
    }
}
