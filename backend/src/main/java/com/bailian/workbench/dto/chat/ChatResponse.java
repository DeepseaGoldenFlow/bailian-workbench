package com.bailian.workbench.dto.chat;

import java.util.List;

public record ChatResponse(
        String id,
        String model,
        List<ChatChoice> choices,
        ChatUsage usage
) {
    public record ChatUsage(
            Integer promptTokens,
            Integer completionTokens,
            Integer totalTokens
    ) {}
}
