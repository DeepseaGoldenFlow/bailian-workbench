package com.bailian.workbench.dto.generic;

import java.util.Map;

public record GenericRequest(
        String model,
        String apiPath,
        Map<String, Object> input,
        Map<String, Object> parameters,
        Boolean async
) {}
