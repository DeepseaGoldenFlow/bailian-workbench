package com.bailian.workbench.dto.embedding;

public record EmbeddingRequest(
        String model,
        Input input,
        EmbeddingParameters parameters
) {
    public record Input(
            String[] texts
    ) {}
    public record EmbeddingParameters(
            String textType
    ) {}
}
