package com.bailian.workbench.service;

import com.bailian.workbench.dto.chat.ChatResponse;
import com.bailian.workbench.dto.vision.VisionRequest;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
public class VisionService {

    private final WebClient webClient;

    public VisionService(WebClient webClient) {
        this.webClient = webClient;
    }

    /**
     * Synchronous vision (multimodal) chat completion.
     * Uses compatible-mode endpoint with qwen-vl-max model.
     */
    public Mono<ChatResponse> chat(VisionRequest request) {
        return webClient.post()
                .uri("/compatible-mode/v1/chat/completions")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(ChatResponse.class);
    }

    /**
     * Streaming vision chat completion via SSE.
     */
    public Flux<String> chatStream(VisionRequest request) {
        // Build a streaming request
        Object streamRequest = buildStreamRequest(request);

        return webClient.post()
                .uri("/compatible-mode/v1/chat/completions")
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.TEXT_EVENT_STREAM)
                .bodyValue(streamRequest)
                .retrieve()
                .bodyToFlux(String.class);
    }

    private Object buildStreamRequest(VisionRequest request) {
        // Use a Map-based approach to add stream=true to the request
        java.util.Map<String, Object> map = new java.util.LinkedHashMap<>();
        map.put("model", request.model() != null ? request.model() : "qwen-vl-max");
        map.put("messages", request.messages());
        map.put("stream", true);
        if (request.parameters() != null) {
            map.putAll(request.parameters());
        }
        return map;
    }
}
