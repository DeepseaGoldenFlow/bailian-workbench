package com.bailian.workbench.dto.common;

public record ErrorResponse(
        Integer code,
        String message,
        String requestId
) {}
