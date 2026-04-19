package com.bailian.workbench.controller;

import com.bailian.workbench.dto.common.ErrorResponse;
import com.bailian.workbench.dto.config.ModelConfig;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/config")
public class ModelConfigController {

    @GetMapping("/models")
    public ResponseEntity<?> getModels() {
        try {
            ModelConfig config = buildModelConfig();
            return ResponseEntity.ok(config);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(new ErrorResponse(500, e.getMessage(), null));
        }
    }

    private ModelConfig buildModelConfig() {
        List<ModelConfig.Category> categories = List.of(
                chatCategory(),
                imageGenCategory(),
                imageEditCategory(),
                tryOnCategory(),
                videoT2vCategory(),
                videoI2vCategory(),
                videoRefCategory(),
                videoEditCategory(),
                digitalHumanCategory(),
                animateMoveCategory(),
                ttsCategory(),
                asrCategory()
        );
        return new ModelConfig(categories);
    }

    private ModelConfig.Category chatCategory() {
        return new ModelConfig.Category("智能对话", "chat", List.of(
                new ModelConfig.Model("qwen3.6-plus", "Qwen 3.6 Plus", "最新旗舰，1M 上下文，深度思考",
                        List.of(
                                new ModelConfig.Param("temperature", "创意程度", "较低时回复精准，较高时更有创造力", "slider", 0.7, 0.0, 2.0, 0.1, null),
                                new ModelConfig.Param("top_p", "回复多样性", "控制模型考虑的候选词范围", "slider", 0.8, 0.0, 1.0, 0.05, null),
                                new ModelConfig.Param("enable_search", "联网搜索", "开启后搜索互联网获取最新信息", "toggle", false, null, null, 0, null)
                        )
                ),
                new ModelConfig.Model("qwen3-max", "Qwen 3 Max", "复杂任务最强推理",
                        List.of(
                                new ModelConfig.Param("temperature", "创意程度", null, "slider", 0.7, 0.0, 2.0, 0.1, null),
                                new ModelConfig.Param("top_p", "回复多样性", null, "slider", 0.8, 0.0, 1.0, 0.05, null)
                        )
                )
        ));
    }

    private ModelConfig.Category imageGenCategory() {
        return new ModelConfig.Category("文生图", "image-gen", List.of(
                new ModelConfig.Model("wan2.7-image-pro", "Wan 2.7 Pro (最强)", "万相 2.7 专业版，最高画质",
                        List.of(
                                new ModelConfig.Param("size", "画面比例", null, "dropdown", "1024*1024", null, null, 0,
                                        List.of("1024*1024", "1620*1080", "1080*1920", "1152*864", "864*1152")),
                                new ModelConfig.Param("n", "生成数量", "一次生成的图片数量", "slider", 1, 1.0, 4.0, 1.0, null),
                                new ModelConfig.Param("negative_prompt", "反向提示词", "描述不想出现的内容", "textarea", "", null, null, 0, null),
                                new ModelConfig.Param("prompt_extend", "提示词优化", "自动优化提示词以提升画质", "toggle", true, null, null, 0, null)
                        )
                ),
                new ModelConfig.Model("wan2.7-image", "Wan 2.7 (快速)", "万相 2.7 标准版，速度更快",
                        List.of(
                                new ModelConfig.Param("size", "画面比例", null, "dropdown", "1024*1024", null, null, 0,
                                        List.of("1024*1024", "1620*1080", "1080*1920", "1152*864")),
                                new ModelConfig.Param("prompt_extend", "提示词优化", null, "toggle", true, null, null, 0, null)
                        )
                )
        ));
    }

    private ModelConfig.Category imageEditCategory() {
        return new ModelConfig.Category("图像编辑", "image-edit", List.of(
                new ModelConfig.Model("wan2.5-i2i-preview", "Wan 2.5 Edit", "风格迁移、编辑、扩图、超分",
                        List.of(
                                new ModelConfig.Param("function", "功能选择", "风格迁移/编辑描述/智能扩图/超分辨率", "dropdown", "style_transfer", null, null, 0,
                                        List.of("style_transfer", "description_edit", "expand", "super_resolution")),
                                new ModelConfig.Param("ref_images", "参考图", "上传 1-4 张参考图片", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("prompt", "编辑指令", "描述你想要的修改效果", "textarea", "", null, null, 0, null),
                                new ModelConfig.Param("mask", "局部重绘蒙版", "仅在局部重绘时需要", "upload", null, null, null, 0, null)
                        )
                )
        ));
    }

    private ModelConfig.Category tryOnCategory() {
        return new ModelConfig.Category("AI 试衣", "try-on", List.of(
                new ModelConfig.Model("aitryon", "AI 试衣", "上传模特和服饰，自动生成试穿效果",
                        List.of(
                                new ModelConfig.Param("model_image", "模特图片", "上传清晰的模特/人物全身照", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("garment_image", "服饰图片", "上传服饰平铺或模特穿着图", "upload", null, null, null, 0, null)
                        )
                )
        ));
    }

    private ModelConfig.Category videoT2vCategory() {
        return new ModelConfig.Category("文生视频", "video-t2v", List.of(
                new ModelConfig.Model("wan2.7-t2v", "Wan 2.7 (最新)", "多镜头叙事，自动配音",
                        List.of(
                                new ModelConfig.Param("resolution", "分辨率", null, "dropdown", "720P", null, null, 0,
                                        List.of("480P", "720P", "1080P")),
                                new ModelConfig.Param("duration", "时长 (秒)", null, "dropdown", "5", null, null, 0,
                                        List.of("5", "10")),
                                new ModelConfig.Param("audio", "自动配音", "根据画面自动生成音效", "toggle", true, null, null, 0, null),
                                new ModelConfig.Param("prompt_extend", "提示词优化", null, "toggle", true, null, null, 0, null)
                        )
                )
        ));
    }

    private ModelConfig.Category videoI2vCategory() {
        return new ModelConfig.Category("图生视频", "video-i2v", List.of(
                new ModelConfig.Model("wan2.7-i2v", "Wan 2.7 (全能)", "首帧/首尾帧/续写三合一",
                        List.of(
                                new ModelConfig.Param("mode", "生成模式", "首帧：图片变视频 / 首尾帧：生成过渡动画 / 续写：延长视频", "dropdown", "first_frame", null, null, 0,
                                        List.of("first_frame", "first_last_frame", "video_continue")),
                                new ModelConfig.Param("first_frame", "首帧图片", "作为视频开头的图片", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("last_frame", "尾帧图片", "作为视频结尾的图片 (仅首尾帧模式)", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("ref_video", "前序视频", "需要续写的原视频 (仅续写模式)", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("resolution", "分辨率", null, "dropdown", "720P", null, null, 0,
                                        List.of("480P", "720P", "1080P")),
                                new ModelConfig.Param("duration", "时长 (秒)", null, "dropdown", "5", null, null, 0,
                                        List.of("5", "10"))
                        )
                )
        ));
    }

    private ModelConfig.Category videoRefCategory() {
        return new ModelConfig.Category("参考生视频", "video-ref", List.of(
                new ModelConfig.Model("wan2.7-r2v", "Wan 2.7 Ref", "保持角色一致性生成",
                        List.of(
                                new ModelConfig.Param("ref_image", "参考图片", "保持画面中的人物/物体特征", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("ref_video", "参考视频", "参考动作和风格", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("resolution", "分辨率", null, "dropdown", "720P", null, null, 0,
                                        List.of("480P", "720P", "1080P")),
                                new ModelConfig.Param("duration", "时长 (秒)", null, "dropdown", "5", null, null, 0,
                                        List.of("5", "10"))
                        )
                )
        ));
    }

    private ModelConfig.Category videoEditCategory() {
        return new ModelConfig.Category("视频编辑", "video-edit", List.of(
                new ModelConfig.Model("wan2.7-videoedit", "Wan 2.7 Edit", "重绘/局部编辑/延展",
                        List.of(
                                new ModelConfig.Param("input_video", "原视频", "需要编辑的原始视频", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("ref_images", "参考图", "用于风格重绘的参考图", "upload-multi", null, null, null, 0, null),
                                new ModelConfig.Param("mask", "局部掩码", "指定修改区域的掩码图", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("resolution", "分辨率", null, "dropdown", "720P", null, null, 0,
                                        List.of("480P", "720P", "1080P")),
                                new ModelConfig.Param("duration", "时长 (秒)", null, "dropdown", "5", null, null, 0,
                                        List.of("5", "10"))
                        )
                )
        ));
    }

    private ModelConfig.Category digitalHumanCategory() {
        return new ModelConfig.Category("数字人", "digital-human", List.of(
                new ModelConfig.Model("wan2.2-s2v", "Wan S2V", "图片+音频生成说话视频",
                        List.of(
                                new ModelConfig.Param("image", "人物图片", "清晰正面人像照片", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("audio", "驱动音频", "驱动口型和表情变化的音频", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("resolution", "分辨率", null, "dropdown", "720P", null, null, 0,
                                        List.of("480P", "720P"))
                        )
                )
        ));
    }

    private ModelConfig.Category animateMoveCategory() {
        return new ModelConfig.Category("图生动作", "animate-move", List.of(
                new ModelConfig.Model("wan2.2-animate-move", "Wan Move", "图片+参考视频生成舞蹈/动作",
                        List.of(
                                new ModelConfig.Param("image", "人物图片", "全身清晰图片，姿势自然", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("ref_video", "动作参考", "包含目标动作/舞蹈的视频", "upload", null, null, null, 0, null),
                                new ModelConfig.Param("mode", "模式", "标准：速度快 / 专业：动作更流畅自然", "dropdown", "wan-std", null, null, 0,
                                        List.of("wan-std", "wan-pro"))
                        )
                )
        ));
    }

    private ModelConfig.Category ttsCategory() {
        return new ModelConfig.Category("语音合成", "tts", List.of(
                new ModelConfig.Model("cosyvoice-v3.5-plus", "CosyVoice V3.5 Plus", "最新最强音质",
                        List.of(
                                new ModelConfig.Param("voice", "音色选择", "不同风格音色", "dropdown", "longwan", null, null, 0,
                                        List.of("longwan", "longcheng", "longhua", "longjing", "longyu")),
                                new ModelConfig.Param("speed", "语速", "0.5 慢 - 2.0 快", "slider", 1.0, 0.5, 2.0, 0.1, null),
                                new ModelConfig.Param("volume", "音量", "0 - 100", "slider", 50, 0, 100, 1, null)
                        )
                ),
                new ModelConfig.Model("qwen-voice-design", "声音设计", "描述生成专属音色",
                        List.of(
                                new ModelConfig.Param("gender", "性别", null, "dropdown", "male", null, null, 0, List.of("male", "female")),
                                new ModelConfig.Param("age", "年龄段", null, "dropdown", "youth", null, null, 0, List.of("child", "youth", "middle", "elderly")),
                                new ModelConfig.Param("desc", "音色描述", "详细描述你想要的声音特质", "textarea", "温柔的", null, null, 0, null)
                        )
                )
        ));
    }

    private ModelConfig.Category asrCategory() {
        return new ModelConfig.Category("语音识别", "asr", List.of(
                new ModelConfig.Model("qwen3-asr-flash", "Qwen3 ASR Flash", "高精度多语种语音识别",
                        List.of(
                                new ModelConfig.Param("audio", "音频文件", "支持 mp3/wav/m4a 等常见格式", "upload", null, null, null, 0, null)
                        )
                )
        ));
    }
}
