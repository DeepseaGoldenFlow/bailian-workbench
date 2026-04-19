package com.bailian.workbench.dto.video;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;

import java.util.List;
import java.util.Map;

@JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
@JsonIgnoreProperties(ignoreUnknown = true)
public record VideoGenRequest(
        String model,
        Input input,
        Map<String, Object> parameters
) {
    /**
     * Flexible input record that supports various video API input formats:
     * - T2V: prompt, audio, prompt_extend
     * - I2V: prompt, first_frame, last_frame, ref_video, multimodal_input
     * - R2V: prompt, ref_image, ref_video
     * - VideoEdit: input_video, ref_images, mask, prompt
     * - S2V (Digital Human): image, audio
     * - Animate Move: image, ref_video
     */
    @JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Input(
            String prompt,
            String imageUrl,
            String audioUrl,
            String firstFrame,
            String lastFrame,
            String refVideo,
            String refImage,
            String inputVideo,
            String mask,
            String image,
            List<String> refImages,
            List<Map<String, Object>> multimodalInput,
            Boolean audio,
            Boolean promptExtend
    ) {}
}
