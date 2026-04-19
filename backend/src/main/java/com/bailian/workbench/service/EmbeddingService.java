package com.bailian.workbench.service;

import com.bailian.workbench.dto.embedding.EmbeddingRequest;
import com.bailian.workbench.dto.embedding.EmbeddingResponse;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

@Service
public class EmbeddingService {

    private final RestClient restClient;

    public EmbeddingService(RestClient restClient) {
        this.restClient = restClient;
    }

    /**
     * Generate text embeddings.
     * Bailian API: POST /api/v1/services/embeddings/text-embedding/text-embedding
     */
    public EmbeddingResponse generate(EmbeddingRequest request) {
        return restClient.post()
                .uri("/api/v1/services/embeddings/text-embedding/text-embedding")
                .body(request)
                .retrieve()
                .body(EmbeddingResponse.class);
    }
}
