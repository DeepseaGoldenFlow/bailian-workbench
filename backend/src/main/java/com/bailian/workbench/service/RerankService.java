package com.bailian.workbench.service;

import com.bailian.workbench.dto.rerank.RerankRequest;
import com.bailian.workbench.dto.rerank.RerankResponse;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

@Service
public class RerankService {

    private final RestClient restClient;

    public RerankService(RestClient restClient) {
        this.restClient = restClient;
    }

    /**
     * Rerank documents by relevance to a query.
     * Bailian API: POST /api/v1/services/rerank/text-rerank/text-rerank
     */
    public RerankResponse rerank(RerankRequest request) {
        return restClient.post()
                .uri("/api/v1/services/rerank/text-rerank/text-rerank")
                .body(request)
                .retrieve()
                .body(RerankResponse.class);
    }
}
