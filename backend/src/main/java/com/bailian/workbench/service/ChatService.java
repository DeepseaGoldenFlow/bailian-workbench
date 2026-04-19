package com.bailian.workbench.service;

import com.bailian.workbench.dto.chat.ChatRequest;
import com.bailian.workbench.dto.chat.ChatResponse;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
public class ChatService {

    private final WebClient webClient;

    public ChatService(WebClient webClient) {
        this.webClient = webClient;
    }

    /**
     * Synchronous chat completion.
     */
    public Mono<ChatResponse> chat(ChatRequest request) {
        return webClient.post()
                .uri("/compatible-mode/v1/chat/completions")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(ChatResponse.class);
    }

    /**
     * Streaming chat completion via SSE.
     */
    public Flux<String> chatStream(ChatRequest request) {
        ChatRequest streamRequest = new ChatRequest(
                request.model(),
                request.messages(),
                request.temperature(),
                request.topP(),
                request.maxTokens(),
                true,
                request.stop()
        );

        return webClient.post()
                .uri("/compatible-mode/v1/chat/completions")
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.TEXT_EVENT_STREAM)
                .bodyValue(streamRequest)
                .retrieve()
                .bodyToFlux(String.class);
    }
}
