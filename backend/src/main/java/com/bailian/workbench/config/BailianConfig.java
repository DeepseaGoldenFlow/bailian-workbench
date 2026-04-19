package com.bailian.workbench.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Data
@Component
@ConfigurationProperties(prefix = "bailian.api")
public class BailianConfig {
    private String baseUrl;
    private String key;
    private Storage storage = new Storage();

    @Data
    public static class Storage {
        private String basePath;
    }
}
