package com.bailian.workbench.dto.chat;

public record ChatChoice(
        Integer index,
        ChatMessage message,
        String finishReason
) {}
