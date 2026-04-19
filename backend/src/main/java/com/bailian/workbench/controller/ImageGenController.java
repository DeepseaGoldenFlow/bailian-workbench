package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.image.ImageGenRequest;
import com.bailian.workbench.dto.image.ImageGenResponse;
import com.bailian.workbench.service.ImageGenService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/services/aigc/text2image")
public class ImageGenController {

    private final ImageGenService imageGenService;

    public ImageGenController(ImageGenService imageGenService) {
        this.imageGenService = imageGenService;
    }

    /**
     * Submit image generation task and return taskId immediately.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/image-synthesis")
    public ResponseEntity<?> submit(@RequestBody ImageGenRequest request) {
        try {
            ImageGenResponse res = imageGenService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
