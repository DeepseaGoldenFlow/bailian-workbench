package com.bailian.workbench.controller;

import com.bailian.workbench.config.BailianConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;

@RestController
@RequestMapping("/api/storage")
public class StorageController {

    private static final Logger log = LoggerFactory.getLogger(StorageController.class);
    private final BailianConfig bailianConfig;

    public StorageController(BailianConfig bailianConfig) {
        this.bailianConfig = bailianConfig;
    }

    /**
     * Serve files from the storage directory.
     * GET /api/storage/{images|videos|audio}/{filename}
     * GET /api/storage/list?dir=images - List files in directory
     */
    @GetMapping({"", "/**"})
    public ResponseEntity<?> serveFile(jakarta.servlet.http.HttpServletRequest request) {
        String requestPath = request.getRequestURI();
        // Remove /api/storage/ prefix
        String relativePath = requestPath.replace("/api/storage/", "");
        // Handle /api/storage/ itself
        if (relativePath.equals(requestPath) || relativePath.isEmpty()) {
            return listDirectory("");
        }

        // If it's a directory listing request (?list)
        if (request.getParameter("list") != null) {
            return listDirectory(relativePath);
        }

        Path basePath = Path.of(bailianConfig.getStorage().getBasePath());
        Path filePath = basePath.resolve(relativePath).normalize();

        // Security check: ensure path is within storage directory
        if (!filePath.startsWith(basePath)) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", "Access denied"));
        }

        if (!Files.exists(filePath)) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", "File not found"));
        }

        if (Files.isDirectory(filePath)) {
            return listDirectory(relativePath);
        }

        try {
            Resource resource = new FileSystemResource(filePath);
            String contentType = Files.probeContentType(filePath);
            MediaType mediaType = contentType != null ? MediaType.parseMediaType(contentType) : MediaType.APPLICATION_OCTET_STREAM;

            return ResponseEntity.ok()
                    .contentType(mediaType)
                    .header(HttpHeaders.CONTENT_DISPOSITION, "inline; filename=\"" + filePath.getFileName() + "\"")
                    .body(resource);
        } catch (Exception e) {
            log.error("Failed to serve file: {}", e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(Map.of("error", "Failed to serve file: " + e.getMessage()));
        }
    }

    private ResponseEntity<?> listDirectory(String relativePath) {
        Path basePath = Path.of(bailianConfig.getStorage().getBasePath());
        Path dirPath = relativePath.isEmpty() ? basePath : basePath.resolve(relativePath).normalize();

        if (!dirPath.startsWith(basePath)) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", "Access denied"));
        }

        if (!Files.exists(dirPath)) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", "Directory not found"));
        }

        if (!Files.isDirectory(dirPath)) {
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("error", "Not a directory"));
        }

        try (Stream<Path> stream = Files.list(dirPath)) {
            var files = stream
                    .filter(Files::isRegularFile)
                    .map(path -> {
                        var m = new java.util.LinkedHashMap<String, Object>();
                        m.put("name", path.getFileName().toString());
                        m.put("size", getFileSize(path));
                        m.put("path", basePath.relativize(path).toString());
                        m.put("modified", getFileModifiedTime(path));
                        return m;
                    })
                    .toList();

            File[] dirFiles = dirPath.toFile().listFiles();
            var dirs = dirFiles != null
                    ? java.util.Arrays.stream(dirFiles)
                            .filter(File::isDirectory)
                            .map(f -> {
                                var m = new java.util.LinkedHashMap<String, Object>();
                                m.put("name", f.getName());
                                m.put("path", basePath.relativize(f.toPath()).toString());
                                return m;
                            })
                            .toList()
                    : java.util.List.<java.util.Map<String, Object>>of();

            return ResponseEntity.ok(Map.of(
                    "files", files,
                    "directories", dirs,
                    "currentPath", relativePath.isEmpty() ? "/" : "/" + relativePath
            ));
        } catch (IOException e) {
            log.error("Failed to list directory: {}", e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(Map.of("error", "Failed to list directory"));
        }
    }

    private long getFileSize(Path path) {
        try {
            return Files.size(path);
        } catch (IOException e) {
            return 0;
        }
    }

    private String getFileModifiedTime(Path path) {
        try {
            return Files.getLastModifiedTime(path).toString();
        } catch (IOException e) {
            return "";
        }
    }
}
