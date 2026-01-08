package detector

import (
    "backend/internal/models"
    "os"
    "path/filepath"
    "runtime"
    "strings"

    "golang.org/x/sys/windows/registry" // for Windows
)


// DetectInstalled перевіряє чи програма встановлена
func DetectInstalled(detection models.Detection) bool {
    if detection.Type == "always" || (len(detection.Paths) == 0 && len(detection.Registry) == 0) {
        return true
    }

    // 1. Перевірка шляхів
    for _, path := range detection.Paths {
        if CheckPathExists(path) {
            return true
        }
    }

    // 2. Перевірка реєстру (тільки Windows)
    if runtime.GOOS == "windows" {
        for _, reg := range detection.Registry {
            if IsOSSupported(reg.OS) && checkRegistry(reg.Key) {
                return true
            }
        }
    }

    return false
}

func ExpandPath(path string) string {
    expandedPath := os.ExpandEnv(path)

    if runtime.GOOS == "windows" {
        replacer := strings.NewReplacer(
            "%AppData%", os.Getenv("AppData"),
            "%LocalAppData%", os.Getenv("LocalAppData"),
            "%ProgramFiles%", os.Getenv("ProgramFiles"),
            "%ProgramFiles(x86)%", os.Getenv("ProgramFiles(x86)"),
            "%UserProfile%", os.Getenv("UserProfile"),
            "%SystemRoot%", os.Getenv("SystemRoot"),
            "%TEMP%", os.Getenv("TEMP"),
            )
        expandedPath = replacer.Replace(expandedPath)
    }

    return expandedPath
}

// CheckPathExists перевіряє чи існує шлях
func CheckPathExists(path string) bool {
    // Розширити змінні оточення
    expanded := ExpandPath(path)

    // Glob для wildcards
    if strings.Contains(expanded, "*") {
        matches, _ := filepath.Glob(expanded)
        return len(matches) > 0
    }

    // Простий stat
    _, err := os.Stat(expanded)
    return err == nil
}

// checkRegistry перевіряє Windows реєстр
func checkRegistry(keyPath string) bool {
    if runtime.GOOS != "windows" {
        return false
    }

    // Розібрати ключ: HKLM\SOFTWARE\...
    parts := strings.SplitN(keyPath, "\\", 2)
    if len(parts) != 2 {
        return false
    }

    var rootKey registry.Key
    switch parts[0] {
    case "HKLM":
        rootKey = registry.LOCAL_MACHINE
    case "HKCU":
        rootKey = registry.CURRENT_USER
    default:
        return false
    }

    // Відкрити ключ
    k, err := registry.OpenKey(rootKey, parts[1], registry.QUERY_VALUE)
    if err != nil {
        return false
    }
    defer k.Close()

    return true
}

// IsOSSupported Is OS supported for this operation
func IsOSSupported(osList []string) bool {
    if len(osList) == 0 {
        return true
    }

    currentOS := runtime.GOOS
    for _, os := range osList {
        if strings.EqualFold(os, currentOS) {
            return true
        }
    }
    return false
}
