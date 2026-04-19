package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.embedding.EmbeddingRequest;
import com.bailian.workbench.dto.embedding.EmbeddingResponse;
import com.bailian.workbench.service.EmbeddingService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/embedding")
public class EmbeddingController {

    private final EmbeddingService embeddingService;

    public EmbeddingController(EmbeddingService embeddingService) {
        this.embeddingService = embeddingService;
    }

    /**
     * Generate text embeddings.
     * POST /api/embedding/generate
     */
    @PostMapping("/generate")
    public ResponseEntity<?> generate(@RequestBody EmbeddingRequest request) {
        try {
            EmbeddingResponse res = embeddingService.generate(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
