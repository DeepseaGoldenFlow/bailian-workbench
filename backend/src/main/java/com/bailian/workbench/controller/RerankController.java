package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.rerank.RerankRequest;
import com.bailian.workbench.dto.rerank.RerankResponse;
import com.bailian.workbench.service.RerankService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/rerank")
public class RerankController {

    private final RerankService rerankService;

    public RerankController(RerankService rerankService) {
        this.rerankService = rerankService;
    }

    /**
     * Rerank documents by relevance to a query.
     * POST /api/rerank/generate
     */
    @PostMapping("/generate")
    public ResponseEntity<?> generate(@RequestBody RerankRequest request) {
        try {
            RerankResponse res = rerankService.rerank(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
