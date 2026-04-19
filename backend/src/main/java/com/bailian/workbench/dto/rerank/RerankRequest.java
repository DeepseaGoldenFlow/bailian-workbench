package com.bailian.workbench.dto.rerank;

public record RerankRequest(
        String model,
        Input input
) {
    public record Input(
            String query,
            String[] documents
    ) {}
}
