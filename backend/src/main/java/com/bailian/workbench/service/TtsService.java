package com.bailian.workbench.service;

import com.bailian.workbench.config.BailianConfig;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;
import java.util.Map;

@Service
public class TtsService {

    private final WebClient webClient;
    private final String apiKey;

    public TtsService(BailianConfig bailianConfig, WebClient webClient) {
        this.webClient = webClient;
        // 确保 API Key 正确
        String key = bailianConfig.getKey();
        if (key == null || key.contains("your-api-key")) {
            this.apiKey = System.getenv("BAILIAN_API_KEY");
        } else {
            this.apiKey = key;
        }
    }

    public Mono<byte[]> generateSpeech(Map<String, Object> request) {
        return webClient.post()
                .uri("https://dashscope.aliyuncs.com/compatible-mode/v1/audio/speech")
                .header("Authorization", "Bearer " + this.apiKey)
                .header("Content-Type", MediaType.APPLICATION_JSON_VALUE)
                .bodyValue(request)
                .retrieve()
                .bodyToMono(byte[].class)
                .doOnError(e -> System.err.println("[TTS] DashScope Error: " + e.getMessage()));
    }
}
