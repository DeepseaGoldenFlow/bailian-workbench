package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.video.VideoGenRequest;
import com.bailian.workbench.dto.video.VideoGenResponse;
import com.bailian.workbench.service.VideoEditService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/video")
public class VideoEditController {

    private final VideoEditService videoEditService;

    public VideoEditController(VideoEditService videoEditService) {
        this.videoEditService = videoEditService;
    }

    /**
     * Submit video editing task (Wan 2.7 VideoEdit) and return taskId immediately.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/edit")
    public ResponseEntity<?> submit(@RequestBody VideoGenRequest request) {
        try {
            VideoGenResponse res = videoEditService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
