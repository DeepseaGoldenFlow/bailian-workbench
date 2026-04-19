package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.generic.GenericRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.*;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestClient;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.Map;

@RestController
@RequestMapping("/api/proxy")
public class ProxyController {

    private static final Logger log = LoggerFactory.getLogger(ProxyController.class);

    private final RestClient restClient;
    private final WebClient webClient;
    private final String baseUrl;

    public ProxyController(RestClient restClient, WebClient webClient,
                           @org.springframework.beans.factory.annotation.Value("${bailian.api.base-url}") String baseUrl) {
        this.restClient = restClient;
        this.webClient = webClient;
        this.baseUrl = baseUrl;
    }

    /**
     * Generic proxy endpoint - forwards requests to Bailian API.
     * POST /api/proxy/{apiPath} with JSON body
     */
    @PostMapping("/{apiPath}")
    public ResponseEntity<?> proxyPost(@PathVariable String apiPath,
                                       @RequestBody(required = false) Map<String, Object> body,
                                       @RequestHeader(value = "X-DashScope-Async", required = false) String asyncHeader,
                                       @RequestParam(required = false) Boolean stream) {
        try {
            log.info("Proxy POST to: {}", apiPath);

            WebClient.RequestBodySpec spec = webClient.post()
                    .uri("/" + apiPath)
                    .contentType(MediaType.APPLICATION_JSON);

            if (asyncHeader != null) {
                spec.header("X-DashScope-Async", asyncHeader);
            }

            if (Boolean.TRUE.equals(stream)) {
                spec.accept(MediaType.TEXT_EVENT_STREAM);
            }

            Object requestBody = body != null ? body : Map.of();

            if (Boolean.TRUE.equals(stream)) {
                return ResponseEntity.ok()
                        .contentType(MediaType.TEXT_EVENT_STREAM)
                        .body(spec.bodyValue(requestBody).retrieve().bodyToFlux(String.class));
            }

            String response = spec.bodyValue(requestBody)
                    .retrieve()
                    .bodyToMono(String.class)
                    .block();

            return ResponseEntity.ok()
                    .contentType(MediaType.APPLICATION_JSON)
                    .body(response);
        } catch (Exception e) {
            log.error("Proxy POST failed: {}", e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, "Proxy error: " + e.getMessage(), null));
        }
    }

    /**
     * Generic proxy endpoint - forwards GET requests to Bailian API.
     * GET /api/proxy/{apiPath}
     */
    @GetMapping("/{apiPath}")
    public ResponseEntity<?> proxyGet(@PathVariable String apiPath,
                                      @RequestParam Map<String, String> params) {
        try {
            log.info("Proxy GET to: {}", apiPath);

            // Build URI with query params
            String uri = "/" + apiPath;
            if (!params.isEmpty()) {
                StringBuilder sb = new StringBuilder(uri);
                sb.append("?");
                params.forEach((k, v) -> sb.append(k).append("=").append(v).append("&"));
                uri = sb.substring(0, sb.length() - 1);
            }

            String response = restClient.get()
                    .uri(uri)
                    .retrieve()
                    .body(String.class);

            return ResponseEntity.ok()
                    .contentType(MediaType.APPLICATION_JSON)
                    .body(response);
        } catch (Exception e) {
            log.error("Proxy GET failed: {}", e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, "Proxy error: " + e.getMessage(), null));
        }
    }
}
