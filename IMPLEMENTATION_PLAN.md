# Bailian All-in-One Workbench — Implementation Plan

## Overview
A unified workbench for Alibaba Cloud Bailian (DashScope) AI services:
- **Chat** — Text & Vision models with streaming SSE
- **Image Generation** — Wanx models (async task + polling)
- **Video Generation** — Wanx video models (async task + polling)

## Architecture

### Backend (Spring Boot 3, Java 21)
- **RestClient** (Spring 6.x native) for all Bailian API calls
- **No Spring AI SDK** — raw HTTP with full control
- **SSE Streaming** for chat via `ResponseBodyExtractor`
- **Scheduled Polling** for async media generation tasks
- **Local File Storage** with configurable disk mapping

### Frontend (Vue 3 + Vite + Element Plus)
- **Tabbed Interface** — Chat / Image Gen / Video Gen
- **Real-time streaming** chat with markdown rendering
- **Progress tracking** for async media tasks
- **Gallery view** for generated media

## Bailian API Patterns

### Chat Completion
```
POST https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
Headers: Authorization: Bearer <API_KEY>
Body: { model, messages[], stream: true }
Response: SSE stream of data: {...} chunks
```

### Image Generation (Async)
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis
Headers: Authorization: Bearer <API_KEY>, X-DashScope-Async: enable
Body: { model, input: { prompt }, parameters: { n, size } }
Response: { output: { task_id, task_status } }

GET /api/v1/tasks/{task_id}
Response: { output: { task_status, results: [{ url }] } }
```

### Video Generation (Async)
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis
Headers: Authorization: Bearer <API_KEY>, X-DashScope-Async: enable
Body: { model, input: { prompt, img_url? }, parameters: {...} }
Response: { output: { task_id, task_status } }

GET /api/v1/tasks/{task_id}
Response: { output: { task_status, results: [{ url }] } }
```

## File Structure
```
bailian-workbench/
├── IMPLEMENTATION_PLAN.md
├── docker-compose.yml
├── backend/
│   ├── Dockerfile
│   ├── pom.xml
│   └── src/main/
│       ├── java/com/bailian/
│       │   ├── BailianApplication.java
│       │   ├── config/
│       │   │   ├── BailianApiConfig.java      (RestClient bean)
│       │   │   ├── FileStorageConfig.java     (storage paths)
│       │   │   └── WebConfig.java             (CORS)
│       │   ├── controller/
│       │   │   ├── ChatController.java        (SSE endpoint)
│       │   │   ├── ImageGenController.java    (submit + status)
│       │   │   └── VideoGenController.java    (submit + status)
│       │   ├── service/
│       │   │   ├── ChatService.java           (streaming chat)
│       │   │   ├── ImageGenService.java       (async image tasks)
│       │   │   └── VideoGenService.java       (async video tasks)
│       │   ├── model/
│       │   │   ├── ChatMessage.java
│       │   │   ├── ChatRequest.java
│       │   │   ├── GenTask.java               (task state tracking)
│       │   │   └── ApiResponse.java
│       │   └── async/
│       │       └── TaskPollingService.java    (scheduled poller)
│       └── resources/
│           └── application.yml
└── frontend/
    ├── Dockerfile
    ├── package.json
    ├── vite.config.js
    ├── index.html
    └── src/
        ├── main.js
        ├── App.vue
        ├── api/index.js
        └── components/
            ├── ChatPanel.vue
            ├── ImageGenPanel.vue
            └── VideoGenPanel.vue
```

## Implementation Phases
1. ✅ **Scaffold** — Project structure, build configs, Docker
2. 🔲 **Backend Core** — RestClient setup, config classes
3. 🔲 **Chat Service** — SSE streaming, vision support
4. 🔲 **Image Service** — Async submit, polling, download
5. 🔲 **Video Service** — Async submit, polling, download
6. 🔲 **Frontend** — Vue 3 app, components, API integration
7. 🔲 **Docker** — Compose orchestration, multi-stage builds
8. 🔲 **Test & Deploy** — End-to-end validation
