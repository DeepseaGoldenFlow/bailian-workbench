package com.bailian.workbench.dto.chat;

import java.util.List;

public record ChatRequest(
        String model,
        List<ChatMessage> messages,
        Double temperature,
        Double topP,
        Integer maxTokens,
        Boolean stream,
        List<String> stop
) {
    public ChatRequest {
        // stream defaults to false
    }

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private String model;
        private List<ChatMessage> messages;
        private Double temperature;
        private Double topP;
        private Integer maxTokens;
        private Boolean stream;
        private List<String> stop;

        public Builder model(String model) { this.model = model; return this; }
        public Builder messages(List<ChatMessage> messages) { this.messages = messages; return this; }
        public Builder temperature(Double temperature) { this.temperature = temperature; return this; }
        public Builder topP(Double topP) { this.topP = topP; return this; }
        public Builder maxTokens(Integer maxTokens) { this.maxTokens = maxTokens; return this; }
        public Builder stream(Boolean stream) { this.stream = stream; return this; }
        public Builder stop(List<String> stop) { this.stop = stop; return this; }

        public ChatRequest build() {
            return new ChatRequest(model, messages, temperature, topP, maxTokens, stream, stop);
        }
    }
}
