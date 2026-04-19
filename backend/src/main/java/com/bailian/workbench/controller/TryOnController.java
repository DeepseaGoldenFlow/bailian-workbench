package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.tryon.TryOnRequest;
import com.bailian.workbench.dto.tryon.TryOnResponse;
import com.bailian.workbench.service.TryOnService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/image")
public class TryOnController {

    private final TryOnService tryOnService;

    public TryOnController(TryOnService tryOnService) {
        this.tryOnService = tryOnService;
    }

    /**
     * Submit AI Try-On task and return taskId immediately.
     * Frontend should poll /api/v1/tasks/{taskId} for status.
     */
    @PostMapping("/tryon")
    public ResponseEntity<?> submit(@RequestBody TryOnRequest request) {
        try {
            TryOnResponse res = tryOnService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
