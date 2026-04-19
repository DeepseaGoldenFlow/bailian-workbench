package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.video.VideoGenRequest;
import com.bailian.workbench.dto.video.VideoGenResponse;
import com.bailian.workbench.service.VideoI2vService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/video")
public class VideoI2vController {

    private final VideoI2vService videoI2vService;

    public VideoI2vController(VideoI2vService videoI2vService) {
        this.videoI2vService = videoI2vService;
    }

    /**
     * Submit image-to-video task (Wan 2.7 I2V) and return taskId immediately.
     * Supports first_frame, first_last_frame, and video_continue modes.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/i2v")
    public ResponseEntity<?> submit(@RequestBody VideoGenRequest request) {
        try {
            VideoGenResponse res = videoI2vService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
