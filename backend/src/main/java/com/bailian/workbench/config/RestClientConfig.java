package com.bailian.workbench.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.http.client.JdkClientHttpRequestFactory;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestClient;
import org.springframework.web.reactive.function.client.ExchangeStrategies;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.time.Duration;

@Configuration
public class RestClientConfig {

    private final BailianConfig bailianConfig;
    private final ObjectMapper objectMapper;

    public RestClientConfig(BailianConfig bailianConfig, ObjectMapper objectMapper) {
        this.bailianConfig = bailianConfig;
        this.objectMapper = objectMapper;
    }

    @Bean
    public RestClient restClient() {
        MappingJackson2HttpMessageConverter converter = new MappingJackson2HttpMessageConverter(objectMapper);
        converter.setSupportedMediaTypes(java.util.List.of(MediaType.APPLICATION_JSON, MediaType.ALL));

        JdkClientHttpRequestFactory requestFactory = new JdkClientHttpRequestFactory(java.net.http.HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(30))
                .build());
        requestFactory.setReadTimeout(Duration.ofMinutes(5));

        return RestClient.builder()
                .baseUrl(bailianConfig.getBaseUrl())
                .defaultHeader("Authorization", "Bearer " + bailianConfig.getKey())
                .defaultHeader("Content-Type", "application/json")
                .messageConverters(converters -> converters.add(converter))
                .requestFactory(requestFactory)
                .build();
    }

    @Bean
    public WebClient webClient() {
        ExchangeStrategies strategies = ExchangeStrategies.builder()
                .codecs(configurer -> configurer.defaultCodecs().jackson2JsonEncoder(
                        new org.springframework.http.codec.json.Jackson2JsonEncoder(objectMapper)))
                .codecs(configurer -> configurer.defaultCodecs().jackson2JsonDecoder(
                        new org.springframework.http.codec.json.Jackson2JsonDecoder(objectMapper)))
                .build();

        return WebClient.builder()
                .baseUrl(bailianConfig.getBaseUrl())
                .defaultHeader("Authorization", "Bearer " + bailianConfig.getKey())
                .defaultHeader("Content-Type", "application/json")
                .exchangeStrategies(strategies)
                .codecs(configurer -> configurer.defaultCodecs().maxInMemorySize(50 * 1024 * 1024))
                .build();
    }
}
