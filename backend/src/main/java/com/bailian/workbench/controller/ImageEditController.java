package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.imageedit.ImageEditRequest;
import com.bailian.workbench.dto.imageedit.ImageEditResponse;
import com.bailian.workbench.service.ImageEditService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/image")
public class ImageEditController {

    private final ImageEditService imageEditService;

    public ImageEditController(ImageEditService imageEditService) {
        this.imageEditService = imageEditService;
    }

    /**
     * Submit image editing task and return taskId immediately.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/edit")
    public ResponseEntity<?> submit(@RequestBody ImageEditRequest request) {
        try {
            ImageEditResponse res = imageEditService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
