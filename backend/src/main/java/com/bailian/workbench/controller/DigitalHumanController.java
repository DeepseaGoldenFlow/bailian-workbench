package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.video.VideoGenRequest;
import com.bailian.workbench.dto.video.VideoGenResponse;
import com.bailian.workbench.service.DigitalHumanService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/video")
public class DigitalHumanController {

    private final DigitalHumanService digitalHumanService;

    public DigitalHumanController(DigitalHumanService digitalHumanService) {
        this.digitalHumanService = digitalHumanService;
    }

    /**
     * Submit digital human (speech-to-video) task and return taskId immediately.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/digital-human")
    public ResponseEntity<?> submit(@RequestBody VideoGenRequest request) {
        try {
            VideoGenResponse res = digitalHumanService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
