package com.bailian.workbench.controller;

import com.bailian.workbench.dto.asr.AsrRequest;
import com.bailian.workbench.dto.asr.AsrResponse;
import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.service.AsrService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/asr")
public class AsrController {

    private final AsrService asrService;

    public AsrController(AsrService asrService) {
        this.asrService = asrService;
    }

    /**
     * Upload audio file and transcribe.
     * POST /api/asr/transcribe with multipart form data.
     */
    @PostMapping(value = "/transcribe", consumes = "multipart/form-data")
    public ResponseEntity<?> transcribe(
            @RequestParam("file") MultipartFile file,
            @RequestParam(value = "model", defaultValue = "paraformer-v2") String model,
            @RequestParam(value = "sample_rate", required = false) String sampleRate) {
        try {
            AsrResponse res = asrService.transcribe(file, model, sampleRate);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }

    /**
     * Submit ASR task with audio URL directly.
     */
    @PostMapping("/submit")
    public ResponseEntity<?> submit(@RequestBody AsrRequest request) {
        try {
            AsrResponse res = asrService.submitTask(request);
            return ResponseEntity.ok(res);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }
}
