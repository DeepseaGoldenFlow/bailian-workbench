package com.bailian.workbench.controller;

import com.bailian.workbench.service.TtsService;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;
import java.util.Map;

@RestController
@RequestMapping("/api/compatible-mode/v1")
public class TtsController {

    private final TtsService ttsService;

    public TtsController(TtsService ttsService) {
        this.ttsService = ttsService;
    }

    @PostMapping("/audio/speech")
    public Mono<ResponseEntity<byte[]>> generateAudio(@RequestBody Map<String, Object> request) {
        return ttsService.generateSpeech(request)
                .map(bytes -> ResponseEntity.ok()
                        .contentType(MediaType.parseMediaType("audio/mpeg"))
                        .header("Content-Disposition", "attachment; filename=speech.mp3")
                        .body(bytes))
                .onErrorResume(e -> {
                    System.err.println("[TTS] Error: " + e.getMessage());
                    return Mono.just(ResponseEntity.internalServerError()
                            .body(("TTS Failed: " + e.getMessage()).getBytes()));
                });
    }
}
