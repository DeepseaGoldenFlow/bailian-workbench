package com.bailian.workbench.dto.vision;

import java.util.List;
import java.util.Map;

public record VisionRequest(
        String model,
        List<ChatMessage> messages,
        Map<String, Object> parameters
) {
    public record ChatMessage(
            String role,
            List<ContentItem> content
    ) {}
    public record ContentItem(
            String type,
            String text,
            ImageUrl imageUrl
    ) {
        public record ImageUrl(String url) {}
    }
}
