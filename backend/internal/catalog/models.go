package catalog

func BuildCatalog() *Catalog {
	f := func(v float64) *float64 { return &v }
	input := func(name, label, typ string, required bool) ParamDef {
		return ParamDef{Name: name, Label: label, Type: typ, Scope: ScopeInput, Required: required}
	}
	param := func(name, label, typ string, def any) ParamDef {
		return ParamDef{Name: name, Label: label, Type: typ, Scope: ScopeParameters, Default: def}
	}
	imageMessageParams := func(family string) []ParamDef {
		maxN, defaultSize := 6.0, "2048*2048"
		if family == "wan27" {
			maxN, defaultSize = 12, "2K"
		}
		p := []ParamDef{
			input("prompt", "提示词", TypeString, true),
			{Name: "images", Label: "参考图片", Type: TypeMediaList, Scope: ScopeInput, Description: "每行一个公网或 OSS 图片 URL；留空即文生图", Placeholder: "https://example.com/reference.png"},
			{Name: "size", Label: "输出尺寸", Type: TypeString, Scope: ScopeParameters, Default: defaultSize, Description: "支持 1K/2K/4K 或 宽*高，具体范围由模型决定"},
			{Name: "n", Label: "生成数量", Type: TypeInt, Scope: ScopeParameters, Default: 1, Min: f(1), Max: f(maxN), Step: f(1)},
			{Name: "seed", Label: "随机种子", Type: TypeInt, Scope: ScopeParameters, Min: f(0), Max: f(2147483647), Step: f(1)},
			param("watermark", "添加 AI 水印", TypeBool, false),
			param("prompt_extend", "智能扩写提示词", TypeBool, true),
		}
		if family == "wan27" {
			p = append(p,
				param("enable_sequential", "组图生成", TypeBool, false),
				param("thinking_mode", "思考模式", TypeBool, true),
				ParamDef{Name: "color_palette", Label: "颜色主题 JSON", Type: TypeJSON, Scope: ScopeParameters, Description: "3-10 个 {hex, ratio} 对象；仅组图关闭时可用", Placeholder: `[{"hex":"#C2D1E6","ratio":"100%"}]`},
			)
		}
		return p
	}
	imageSimpleParams := func() []ParamDef {
		return []ParamDef{
			input("prompt", "提示词", TypeString, true), input("negative_prompt", "反向提示词", TypeString, false),
			{Name: "ref_img", Label: "参考图片 URL", Type: TypeMedia, Scope: ScopeInput},
			{Name: "size", Label: "输出尺寸", Type: TypeString, Scope: ScopeParameters, Default: "1024*1024", Placeholder: "1024*1024"},
			{Name: "n", Label: "生成数量", Type: TypeInt, Scope: ScopeParameters, Default: 1, Min: f(1), Max: f(4), Step: f(1)},
			{Name: "seed", Label: "随机种子", Type: TypeInt, Scope: ScopeParameters, Min: f(0), Max: f(2147483647), Step: f(1)},
			param("prompt_extend", "智能扩写提示词", TypeBool, true), param("watermark", "添加 AI 水印", TypeBool, false),
		}
	}
	videoParams := func(maxDuration float64, allInOne bool) []ParamDef {
		return []ParamDef{
			{Name: "prompt", Label: "提示词 / 编辑指令", Type: TypeString, Scope: ScopeInput, Required: !allInOne, Description: "描述画面、动作、运镜、对白与声音"},
			{Name: "media", Label: "输入媒体", Type: TypeMediaList, Scope: ScopeInput, Description: "每行：类型 | URL。类型可用 first_frame、last_frame、reference_image、reference_video、audio、file、link", Placeholder: "first_frame | https://example.com/start.png"},
			{Name: "resolution", Label: "分辨率", Type: TypeSelect, Scope: ScopeParameters, Default: "1080P", Options: []Option{{"1080P", "1080P"}, {"720P", "720P"}, {"480P", "480P"}}},
			{Name: "ratio", Label: "宽高比", Type: TypeSelect, Scope: ScopeParameters, Default: "adaptive", Options: []Option{{"自适应", "adaptive"}, {"16:9", "16:9"}, {"9:16", "9:16"}, {"1:1", "1:1"}, {"4:3", "4:3"}, {"3:4", "3:4"}}},
			{Name: "duration", Label: "时长（秒）", Type: TypeInt, Scope: ScopeParameters, Default: 5, Min: f(2), Max: f(maxDuration), Step: f(1)},
			param("audio", "生成音频", TypeBool, true), param("prompt_extend", "智能扩写提示词", TypeBool, true),
			{Name: "seed", Label: "随机种子", Type: TypeInt, Scope: ScopeParameters, Min: f(0), Max: f(2147483647), Step: f(1)},
			param("watermark", "添加 AI 水印", TypeBool, false),
		}
	}
	happyHorseParams := func() []ParamDef {
		return []ParamDef{
			input("prompt", "提示词 / 编辑指令", TypeString, true),
			{Name: "media", Label: "输入媒体", Type: TypeMediaList, Scope: ScopeInput, Description: "每行：类型 | URL。支持 first_frame、reference_image、video", Placeholder: "first_frame | https://example.com/start.png"},
			{Name: "resolution", Label: "分辨率", Type: TypeSelect, Scope: ScopeParameters, Default: "720P", Options: []Option{{"1080P", "1080P"}, {"720P", "720P"}}},
			{Name: "ratio", Label: "宽高比", Type: TypeSelect, Scope: ScopeParameters, Default: "16:9", Options: []Option{{"16:9", "16:9"}, {"9:16", "9:16"}, {"1:1", "1:1"}, {"4:3", "4:3"}, {"3:4", "3:4"}}},
			{Name: "duration", Label: "时长（秒）", Type: TypeInt, Scope: ScopeParameters, Default: 5, Min: f(3), Max: f(15), Step: f(1)},
			{Name: "seed", Label: "随机种子", Type: TypeInt, Scope: ScopeParameters, Min: f(0), Max: f(2147483647), Step: f(1)},
			param("watermark", "添加 AI 水印", TypeBool, true),
			{Name: "audio_setting", Label: "声音控制", Type: TypeSelect, Scope: ScopeParameters, Default: "auto", Options: []Option{{"自动", "auto"}, {"保留原声", "origin"}}},
		}
	}

	const imageMessagesEndpoint = "/api/v1/services/aigc/multimodal-generation/generation"
	const imageAsyncMessagesEndpoint = "/api/v1/services/aigc/image-generation/generation"
	const imageSimpleEndpoint = "/api/v1/services/aigc/text2image/image-synthesis"
	const videoEndpoint = "/api/v1/services/aigc/video-generation/video-synthesis"
	models := []ModelEntry{
		{ID: "qwen-image-3.0-pro", Name: "Qwen Image 3.0 Pro", Category: CatImage, Endpoint: imageMessagesEndpoint, Payload: "messages", Description: "复杂版面、小字与多语言文字渲染旗舰模型", Parameters: imageMessageParams("qwen")},
		{ID: "qwen-image-3.0", Name: "Qwen Image 3.0", Category: CatImage, Endpoint: imageMessagesEndpoint, Payload: "messages", Description: "Qwen Image 3.0 高速版", Parameters: imageMessageParams("qwen")},
		{ID: "qwen-image-2.0-pro", Name: "Qwen Image 2.0 Pro", Category: CatImage, Endpoint: imageMessagesEndpoint, Payload: "messages", Description: "支持文生图与多图编辑", Parameters: imageMessageParams("qwen")},
		{ID: "qwen-image-2.0", Name: "Qwen Image 2.0", Category: CatImage, Endpoint: imageMessagesEndpoint, Payload: "messages", Description: "支持文生图与多图编辑的快速版", Parameters: imageMessageParams("qwen")},
		{ID: "wan2.7-image-pro", Name: "Wan 2.7 Image Pro", Category: CatImage, Endpoint: imageAsyncMessagesEndpoint, Async: true, Payload: "messages", Description: "最高 4K，支持多图编辑、角色一致性与连续组图", Parameters: imageMessageParams("wan27")},
		{ID: "wan2.7-image", Name: "Wan 2.7 Image", Category: CatImage, Endpoint: imageAsyncMessagesEndpoint, Async: true, Payload: "messages", Description: "最高 2K，支持生成、编辑和连续组图", Parameters: imageMessageParams("wan27")},
		{ID: "wan2.6-image", Name: "Wan 2.6 Image", Category: CatImage, Endpoint: imageAsyncMessagesEndpoint, Async: true, Payload: "messages", Description: "图文混排、多图融合与编辑", Parameters: imageMessageParams("qwen")},
		{ID: "wan2.6-t2i", Name: "Wan 2.6 Text to Image", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "万相 2.6 文生图", Parameters: imageSimpleParams()},
		{ID: "z-image-turbo", Name: "Z-Image Turbo", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "快速低成本写实图像生成", Parameters: imageSimpleParams()},
		{ID: "wan2.5-t2i-preview", Name: "Wan 2.5 Text to Image", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "万相 2.5 文生图", Parameters: imageSimpleParams()},
		{ID: "wan2.2-t2i-plus", Name: "Wan 2.2 T2I Plus", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "万相 2.2 高质量文生图", Parameters: imageSimpleParams()},
		{ID: "wan2.2-t2i-flash", Name: "Wan 2.2 T2I Flash", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "万相 2.2 快速文生图", Parameters: imageSimpleParams()},
		{ID: "wanx2.1-t2i-plus", Name: "Wanx 2.1 T2I Plus", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "兼容旧版万相高质量模型", Parameters: imageSimpleParams()},
		{ID: "wanx2.1-t2i-turbo", Name: "Wanx 2.1 T2I Turbo", Category: CatImage, Endpoint: imageSimpleEndpoint, Async: true, Payload: "standard", Description: "兼容旧版万相极速模型", Parameters: imageSimpleParams()},

		{ID: "wan3.0-video-prime", Name: "Wan 3.0 Video Prime", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "全能高速视频模型：文生、首尾帧、参考、编辑、续写，最长 30 秒", Parameters: videoParams(30, true)},
		{ID: "wan3.0-video", Name: "Wan 3.0 Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "全能视频模型：文生、首尾帧、参考、编辑、续写，最长 30 秒", Parameters: videoParams(30, true)},
		{ID: "wan2.7-t2v", Name: "Wan 2.7 Text to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "有声文生视频与多镜头叙事", Parameters: videoParams(15, false)},
		{ID: "wan2.7-i2v", Name: "Wan 2.7 Image to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "首帧、首尾帧、视频续写与音频驱动", Parameters: videoParams(15, true)},
		{ID: "wan2.7-r2v-2026-06-12", Name: "Wan 2.7 Reference to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "多角色、多图与多视频参考生成", Parameters: videoParams(10, true)},
		{ID: "wan2.7-videoedit", Name: "Wan 2.7 Video Edit", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "视频指令编辑、风格迁移与内容延长", Parameters: videoParams(10, true)},
		{ID: "wan2.6-t2v", Name: "Wan 2.6 Text to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "有声文生视频，最长 15 秒", Parameters: videoParams(15, false)},
		{ID: "wan2.6-i2v", Name: "Wan 2.6 Image to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "有声图生视频，最长 15 秒", Parameters: videoParams(15, true)},
		{ID: "wan2.6-i2v-flash", Name: "Wan 2.6 I2V Flash", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "快速有声图生视频", Parameters: videoParams(15, true)},
		{ID: "wan2.6-r2v", Name: "Wan 2.6 Reference to Video", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "单/多角色参考生视频", Parameters: videoParams(10, true)},
		{ID: "wan2.6-r2v-flash", Name: "Wan 2.6 R2V Flash", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "快速多角色参考生视频", Parameters: videoParams(10, true)},
		{ID: "happyhorse-1.0-t2v", Name: "HappyHorse 文生视频", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "HappyHorse 文本生成视频", Parameters: happyHorseParams()},
		{ID: "happyhorse-1.0-i2v", Name: "HappyHorse 图生视频", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "HappyHorse 首帧图生视频", Parameters: happyHorseParams()},
		{ID: "happyhorse-1.0-r2v", Name: "HappyHorse 参考生视频", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "HappyHorse 最多 9 张参考图生成视频", Parameters: happyHorseParams()},
		{ID: "happyhorse-1.0-video-edit", Name: "HappyHorse 视频编辑", Category: CatVideo, Endpoint: videoEndpoint, Async: true, Payload: "standard", Description: "视频+参考图+指令编辑", Parameters: happyHorseParams()},
	}
	return &Catalog{Models: models}
}
