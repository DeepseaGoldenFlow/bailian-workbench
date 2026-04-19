package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.bailian.workbench.dto.task.TaskResponse;
import com.bailian.workbench.dto.tts.TtsRequest;
import com.bailian.workbench.dto.tts.TtsResponse;
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
public class TtsService {

    private static final Logger log = LoggerFactory.getLogger(TtsService.class);
    private static final int MAX_POLL_ATTEMPTS = 120;
    private static final long POLL_INTERVAL_SECONDS = 5;

    private final RestClient restClient;
    private final WebClient webClient;
    private final TaskService taskService;
    private final BailianConfig bailianConfig;

    public TtsService(RestClient restClient, WebClient webClient, TaskService taskService, BailianConfig bailianConfig) {
        this.restClient = restClient;
        this.webClient = webClient;
        this.taskService = taskService;
        this.bailianConfig = bailianConfig;
    }

    public String generateAudio(TtsRequest request) {
        TtsResponse submitResponse = restClient.post()
                .uri("/api/v1/services/aigc/text2audio/speech-synthesis")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(TtsResponse.class);

        String taskId = submitResponse.output().taskId();
        log.info("TTS task submitted: {}", taskId);

        TtsResponse.Output result = pollTask(taskId);
        return downloadAudio(result);
    }

    public TtsResponse submitTask(TtsRequest request) {
        return restClient.post()
                .uri("/api/v1/services/aigc/text2audio/speech-synthesis")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(TtsResponse.class);
    }

    private TtsResponse.Output pollTask(String taskId) {
        int attempts = 0;
        while (attempts < MAX_POLL_ATTEMPTS) {
            try {
                TaskResponse taskResponse = taskService.getTaskStatus(taskId);
                String status = taskResponse.output().taskStatus();

                if ("SUCCEEDED".equalsIgnoreCase(status)) {
                    return convertOutput(taskResponse.output());
                } else if ("FAILED".equalsIgnoreCase(status)) {
                    throw new RuntimeException("TTS task failed: " + taskResponse.message());
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
        throw new RuntimeException("TTS task timed out");
    }

    private TtsResponse.Output convertOutput(TaskResponse.Output output) {
        String audioUrl = null;
        if (output.results() != null && !output.results().isEmpty()) {
            audioUrl = output.results().get(0).url();
        }
        return new TtsResponse.Output(output.taskId(), output.taskStatus(), audioUrl);
    }

    private String downloadAudio(TtsResponse.Output output) {
        if (output.audioUrl() == null) throw new RuntimeException("No audio result");

        String basePath = bailianConfig.getStorage().getBasePath();
        Path audioDir = Path.of(basePath, "audio");
        try { Files.createDirectories(audioDir); } catch (IOException e) { throw new RuntimeException(e); }

        String filename = UUID.randomUUID() + ".wav";
        Path localPath = audioDir.resolve(filename);
        try {
            byte[] bytes = webClient.get().uri(output.audioUrl())
                    .retrieve()
                    .onStatus(HttpStatusCode::isError, resp -> Mono.error(new RuntimeException("Download failed")))
                    .bodyToMono(byte[].class).block();
            Files.write(localPath, bytes);
            return "/storage/audio/" + filename;
        } catch (Exception e) {
            throw new RuntimeException("Audio download failed: " + e.getMessage());
        }
    }
}
