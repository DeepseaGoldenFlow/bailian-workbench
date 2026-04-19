package com.bailian.workbench.dto.tts;

import java.util.Map;

public record TtsRequest(
        String model,
        Input input,
        Map<String, Object> parameters
) {
    public record Input(
            String text
    ) {}
}
