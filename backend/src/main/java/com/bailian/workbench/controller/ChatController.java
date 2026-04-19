package com.bailian.workbench.controller;

import com.bailian.workbench.dto.chat.ChatRequest;
import com.bailian.workbench.dto.chat.ChatResponse;
import com.bailian.workbench.service.ChatService;
import org.springframework.http.MediaType;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/compatible-mode/v1")
public class ChatController {

    private final ChatService chatService;

    public ChatController(ChatService chatService) {
        this.chatService = chatService;
    }

    /**
     * Smart chat completion endpoint.
     * If stream=true, returns SSE stream; otherwise returns synchronous response.
     */
    @PostMapping(value = "/chat/completions", consumes = MediaType.APPLICATION_JSON_VALUE)
    public Object chatCompletions(@RequestBody ChatRequest request) {
        if (Boolean.TRUE.equals(request.stream())) {
            return chatService.chatStream(request)
                    .map(data -> ServerSentEvent.<String>builder()
                            .event("message")
                            .data(data)
                            .build());
        }
        return chatService.chat(request);
    }
}
