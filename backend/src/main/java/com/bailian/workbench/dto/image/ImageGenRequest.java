package com.bailian.workbench.dto.image;

import java.util.Map;

public record ImageGenRequest(
        String model,
        Input input,
        Map<String, Object> parameters
) {
    public record Input(
            String prompt,
            String negativePrompt
    ) {}
}
