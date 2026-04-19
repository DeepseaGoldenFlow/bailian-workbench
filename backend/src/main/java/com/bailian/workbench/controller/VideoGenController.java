package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.video.VideoGenRequest;
import com.bailian.workbench.dto.video.VideoGenResponse;
import com.bailian.workbench.service.VideoGenService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/services/aigc/video-generation")
public class VideoGenController {

    private final VideoGenService videoGenService;

    public VideoGenController(VideoGenService videoGenService) {
        this.videoGenService = videoGenService;
    }

    @PostMapping("/video-synthesis")
    public ResponseEntity<?> submit(@RequestBody VideoGenRequest request) {
        try {
            VideoGenResponse res = videoGenService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
