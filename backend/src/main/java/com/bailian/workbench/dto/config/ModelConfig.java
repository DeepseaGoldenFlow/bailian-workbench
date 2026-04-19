package com.bailian.workbench.dto.config;

import java.util.List;
import java.util.Map;

public record ModelConfig(
        List<Category> categories
) {
    public record Category(
            String name,
            String icon,
            List<Model> models
    ) {}
    public record Model(
            String id,
            String name,
            String description,
            List<Param> parameters
    ) {}
    public record Param(
            String key,
            String label,
            String description,
            String type,
            Object defaultValue,
            Object min,
            Object max,
            double step,
            List<String> options
    ) {}
}
