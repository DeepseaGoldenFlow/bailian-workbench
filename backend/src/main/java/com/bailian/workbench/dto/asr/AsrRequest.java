package com.bailian.workbench.dto.asr;

import java.util.Map;

public record AsrRequest(
        String model,
        Input input,
        Map<String, Object> parameters
) {
    public record Input(
            String audioUrl
    ) {}
}
