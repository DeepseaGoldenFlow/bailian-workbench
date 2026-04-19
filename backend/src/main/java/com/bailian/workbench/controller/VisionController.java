package com.bailian.workbench.controller;

import com.bailian.workbench.dto.chat.ChatResponse;
import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.vision.VisionRequest;
import com.bailian.workbench.service.VisionService;
import org.springframework.http.MediaType;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/vision")
public class VisionController {

    private final VisionService visionService;

    public VisionController(VisionService visionService) {
        this.visionService = visionService;
    }

    /**
     * Vision (multimodal) chat completion.
     * If stream=true, returns SSE stream; otherwise returns synchronous response.
     * POST /api/vision/chat
     */
    @PostMapping(value = "/chat", consumes = MediaType.APPLICATION_JSON_VALUE)
    public Object chat(@RequestBody VisionRequest request) {
        // Check if streaming is requested via parameters
        boolean stream = request.parameters() != null
                && Boolean.TRUE.equals(request.parameters().get("stream"));

        if (stream) {
            return visionService.chatStream(request)
                    .map(data -> ServerSentEvent.<String>builder()
                            .event("message")
                            .data(data)
                            .build());
        }
        return visionService.chat(request);
    }
}
