package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.service.PptService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping("/api/ppt")
public class PptController {

    private final PptService pptService;

    public PptController(PptService pptService) {
        this.pptService = pptService;
    }

    @PostMapping("/generate")
    public ResponseEntity<?> generate(@RequestBody Map<String, Object> req) {
        try {
            String topic = (String) req.get("topic");
            if (topic == null || topic.trim().isEmpty()) {
                return ResponseEntity.badRequest().body(new ErrorResponse(400, "Topic is required", null));
            }
            int slides = req.get("slide_count") != null ? ((Number) req.get("slide_count")).intValue() : 5;
            String style = req.get("style") != null ? req.get("style").toString() : "商务简约";

            String url = pptService.createPpt(topic, slides, style);
            return ResponseEntity.ok(Map.of("status", "success", "download_url", url));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, "PPT Generation Failed: " + e.getMessage(), null));
        }
    }
}
