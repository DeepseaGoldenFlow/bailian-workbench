# 百炼工作站

一个部署在 NAS 上的阿里云百炼多媒体工作站。图片和视频页面由后端模型目录动态生成表单，不再把模型参数写死在页面里。

## 这版支持什么

- 图片：Qwen Image 3.0/2.0、Wan 2.7/2.6、Z-Image、Wanx 等生成与编辑模型。
- 视频：Wan 3.0、Wan 2.7/2.6 与 HappyHorse，覆盖文生视频、首尾帧、参考生成、编辑和续写。
- 常用参数直接显示为表单，包括尺寸、数量、随机种子、时长、分辨率、宽高比、音频、水印等。
- `input JSON` 和 `parameters JSON` 可覆盖或补充文档中的任意参数，新模型兼容时还可填写“自定义模型 ID”。
- 异步任务会自动轮询，完成后直接显示图片或播放视频。

## 启动

在项目根目录创建 `.env`：

```dotenv
BAILIAN_API_KEY=你的百炼_API_Key
MYSQL_ROOT_PASSWORD=数据库_root_密码
MYSQL_PASSWORD=数据库用户密码
MYSQL_USER=bailian

# 默认是阿里云百炼公网地址；使用兼容网关时可覆盖
DASHSCOPE_BASE_URL=https://dashscope.aliyuncs.com
```

然后启动：

```bash
docker compose up -d --build
```

网页默认端口为 `3800`，后端默认端口为 `8080`，数据库映射端口为 `13307`。

## 输入媒体

图片参考图每行填写一个公网或 OSS URL。视频媒体每行格式为：

```text
first_frame | https://example.com/start.png
last_frame | https://example.com/end.png
audio | https://example.com/voice.mp3
```

Wan 3.0 等模型还可使用 `reference_image`、`reference_video`、`file`、`link` 等类型。若模型文档出现页面尚未列出的字段，直接在高级 JSON 中粘贴完整的 `input` 或 `parameters` 对象即可。

> 媒体 URL 必须能被百炼服务访问。本机文件不能直接作为 URL 提交，应先上传至 OSS 或其他可公开读取的位置。
