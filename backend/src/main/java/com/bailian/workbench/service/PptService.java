package com.bailian.workbench.service;

import com.bailian.workbench.dto.chat.ChatMessage;
import com.bailian.workbench.dto.chat.ChatRequest;
import com.bailian.workbench.dto.chat.ChatResponse;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.poi.xslf.usermodel.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.awt.*;
import java.awt.geom.Rectangle2D;
import java.io.FileOutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Service
public class PptService {

    private static final Logger log = LoggerFactory.getLogger(PptService.class);
    private final ChatService chatService;
    private final ObjectMapper mapper = new ObjectMapper();

    public PptService(ChatService chatService) {
        this.chatService = chatService;
    }

    public String createPpt(String topic, int slideCount, String style) throws Exception {
        // 1. Call AI to generate outline
        log.info("Generating PPT outline for topic: {}", topic);
        
        // Construct prompt for JSON output
        String prompt = String.format(
            "生成一份关于《%s》的PPT大纲，共%d页，风格为%s。\n" +
            "请仅返回一个合法的JSON数组，不要任何其他文字。\n" +
            "格式要求: [{\"title\":\"标题\",\"content\":\"正文要点\"}]\n" +
            "content 字段请用换行符分隔要点。",
            topic, slideCount, style
        );

        // Call ChatService
        ChatRequest request = new ChatRequest(
            "qwen3.6-plus",
            List.of(new ChatMessage("user", prompt)),
            0.7, 0.8, 4000, false, null
        );

        ChatResponse response = chatService.chat(request).block();
        if (response == null || response.choices() == null || response.choices().isEmpty()) {
            throw new RuntimeException("AI 模型调用失败");
        }
        String aiContent = response.choices().get(0).message().content();
        
        // 2. Parse JSON
        // Clean up markdown code blocks if present
        String jsonStr = aiContent.trim();
        Pattern p = Pattern.compile("```(?:json)?\\s*([\\s\\S]*?)```");
        Matcher m = p.matcher(jsonStr);
        if (m.find()) {
            jsonStr = m.group(1).trim();
        }

        List<Map<String, String>> slides;
        try {
            slides = mapper.readValue(jsonStr, new TypeReference<>() {});
        } catch (Exception e) {
            log.error("Failed to parse JSON: {}", jsonStr);
            throw new RuntimeException("AI 返回的大纲格式错误");
        }

        if (slides.isEmpty()) {
            throw new RuntimeException("AI 未生成有效内容");
        }

        // 3. Generate PPTX using Apache POI
        XMLSlideShow ppt = new XMLSlideShow();
        ppt.setPageSize(new Dimension(960, 720)); // 16:9 aspect ratio approx

        // Set default font (Skipped due to API compatibility)
        // XSLFSlideMaster master = ppt.getSlideMasters().get(0);
        // XSLFTitleMaster titleMaster = master.getTitleMaster();
        
        for (Map<String, String> slideData : slides) {
            XSLFSlide slide = ppt.createSlide();
            
            // Add Title
            XSLFTextShape titleShape = slide.createTextBox();
            titleShape.setAnchor(new Rectangle2D.Double(50, 30, 860, 100));
            XSLFTextParagraph titlePara = titleShape.addNewTextParagraph();
            XSLFTextRun titleRun = titlePara.addNewTextRun();
            titleRun.setText(slideData.get("title"));
            titleRun.setFontColor(new Color(51, 51, 51));
            titleRun.setFontSize(32.0);
            titleRun.setBold(true);
            titleRun.setFontFamily("微软雅黑");

            // Add Body
            XSLFTextBox body = slide.createTextBox();
            body.setAnchor(new Rectangle2D.Double(50, 150, 860, 500));
            
            XSLFTextParagraph para = body.addNewTextParagraph();
            XSLFTextRun run = para.addNewTextRun();
            run.setText(slideData.get("content"));
            run.setFontSize(18.0);
            run.setFontColor(new Color(80, 80, 80));
            run.setFontFamily("微软雅黑");
        }

        // 4. Save to Storage
        String filename = "ppt_" + System.currentTimeMillis() + "_" + UUID.randomUUID().toString().substring(0,8) + ".pptx";
        Path storageDir = Path.of("/data/bailian/storage/ppts");
        Files.createDirectories(storageDir);
        Path filePath = storageDir.resolve(filename);

        try (FileOutputStream out = new FileOutputStream(filePath.toFile())) {
            ppt.write(out);
        }

        log.info("PPT saved to: {}", filePath);
        return "/api/storage/ppts/" + filename;
    }
}
