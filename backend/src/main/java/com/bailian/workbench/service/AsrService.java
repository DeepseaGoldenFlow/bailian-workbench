package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.bailian.workbench.dto.asr.AsrRequest;
import com.bailian.workbench.dto.asr.AsrResponse;
import com.bailian.workbench.dto.task.TaskResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service
public class AsrService {

    private static final Logger log = LoggerFactory.getLogger(AsrService.class);
    private static final int MAX_POLL_ATTEMPTS = 120;
    private static final long POLL_INTERVAL_SECONDS = 5;

    private final RestClient restClient;
    private final TaskService taskService;
    private final BailianConfig bailianConfig;

    public AsrService(RestClient restClient, TaskService taskService, BailianConfig bailianConfig) {
        this.restClient = restClient;
        this.taskService = taskService;
        this.bailianConfig = bailianConfig;
    }

    /**
     * Upload audio file to storage, then submit ASR task, poll and return result.
     */
    public AsrResponse transcribe(MultipartFile file, String model, String sampleRate) {
        // 1. Save uploaded file to local storage
        String audioUrl = saveAudioFile(file);

        // 2. Submit ASR task
        AsrRequest request = new AsrRequest(
                model,
                new AsrRequest.Input(audioUrl),
                sampleRate != null ? java.util.Map.of("sample_rate", Integer.parseInt(sampleRate)) : null
        );

        AsrResponse submitResponse = restClient.post()
                .uri("/api/v1/services/aicp/paraformer/async-translation")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(AsrResponse.class);

        String taskId = submitResponse.output().taskId();
        log.info("ASR task submitted: {}", taskId);

        return pollTask(taskId);
    }

    /**
     * Submit ASR task with audio URL directly.
     */
    public AsrResponse submitTask(AsrRequest request) {
        return restClient.post()
                .uri("/api/v1/services/aicp/paraformer/async-translation")
                .header("X-DashScope-Async", "enable")
                .body(request)
                .retrieve()
                .body(AsrResponse.class);
    }

    private String saveAudioFile(MultipartFile file) {
        String basePath = bailianConfig.getStorage().getBasePath();
        Path audioDir = Path.of(basePath, "audio");
        try { Files.createDirectories(audioDir); } catch (IOException e) { throw new RuntimeException(e); }

        String originalFilename = file.getOriginalFilename();
        String extension = originalFilename != null && originalFilename.contains(".")
                ? originalFilename.substring(originalFilename.lastIndexOf("."))
                : ".wav";
        String filename = UUID.randomUUID() + extension;
        Path localPath = audioDir.resolve(filename);

        try {
            file.transferTo(localPath.toFile());
            // Return a URL that the Bailian API can access
            // Since Bailian API needs a public URL, we use the file path as URL
            // In production, you'd upload to OSS or use a public URL
            return "file://" + localPath.toAbsolutePath();
        } catch (IOException e) {
            throw new RuntimeException("Failed to save audio file: " + e.getMessage(), e);
        }
    }

    private AsrResponse pollTask(String taskId) {
        int attempts = 0;
        while (attempts < MAX_POLL_ATTEMPTS) {
            try {
                TaskResponse taskResponse = taskService.getTaskStatus(taskId);
                String status = taskResponse.output().taskStatus();

                if ("SUCCEEDED".equalsIgnoreCase(status)) {
                    return convertOutput(taskResponse.output());
                } else if ("FAILED".equalsIgnoreCase(status)) {
                    throw new RuntimeException("ASR task failed: " + taskResponse.message());
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
        throw new RuntimeException("ASR task timed out");
    }

    private AsrResponse convertOutput(TaskResponse.Output output) {
        // The ASR result text is typically in the result URL or as direct text
        // We need to fetch the actual result from the result URL
        String text = null;
        java.util.List<AsrResponse.Output.Sentence> sentences = null;

        if (output.results() != null && !output.results().isEmpty()) {
            String resultUrl = output.results().get(0).url();
            try {
                // Fetch the result JSON from the result URL
                String resultJson = restClient.get()
                        .uri(resultUrl)
                        .retrieve()
                        .body(String.class);

                // Parse the result - the actual ASR result structure varies
                // Try to extract text from the result
                com.fasterxml.jackson.databind.ObjectMapper mapper = new com.fasterxml.jackson.databind.ObjectMapper();
                com.fasterxml.jackson.databind.JsonNode root = mapper.readTree(resultJson);
                if (root.has("text")) {
                    text = root.get("text").asText();
                }
                if (root.has("sentences")) {
                    sentences = new java.util.ArrayList<>();
                    for (com.fasterxml.jackson.databind.JsonNode s : root.get("sentences")) {
                        String sText = s.has("text") ? s.get("text").asText() : null;
                        Double beginTime = s.has("begin_time") ? s.get("begin_time").asDouble() : null;
                        Double endTime = s.has("end_time") ? s.get("end_time").asDouble() : null;
                        Double confidence = s.has("confidence") ? s.get("confidence").asDouble() : null;
                        sentences.add(new AsrResponse.Output.Sentence(sText, beginTime, endTime, confidence));
                    }
                }
            } catch (Exception e) {
                log.warn("Failed to parse ASR result: {}", e.getMessage());
            }
        }

        return new AsrResponse(
                new AsrResponse.Output(output.taskId(), output.taskStatus(), text, sentences),
                output.taskId() // requestId placeholder
        );
    }
}
