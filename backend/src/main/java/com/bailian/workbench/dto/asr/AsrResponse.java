package com.bailian.workbench.dto.asr;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;
import java.util.List;

@JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
@JsonIgnoreProperties(ignoreUnknown = true)
public record AsrResponse(
        Output output,
        String requestId
) {
    @JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Output(
            String taskId,
            String taskStatus,
            String text,
            List<Sentence> sentences
    ) {
        @JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
        @JsonIgnoreProperties(ignoreUnknown = true)
        public record Sentence(
                String text,
                Double beginTime,
                Double endTime,
                Double confidence
        ) {}
    }
}
