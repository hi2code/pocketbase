package main

import (
	"context"
	"fmt"
	"net/http/httptest"

	"github.com/pocketbase/pocketbase/tools/i18n"
	"golang.org/x/text/language"
)

func main() {
	// Initialize i18n
	if err := i18n.Init(); err != nil {
		panic(err)
	}

	fmt.Println("=== i18n Example ===")
	fmt.Println()

	// Example 1: Basic translation
	fmt.Println("1. Basic Translation:")
	ctxEn := context.Background()
	ctxEn = i18n.SetLanguage(ctxEn, language.English)

	msgEn := i18n.T(ctxEn, "validation_unknown_field")
	fmt.Printf("  English: %s\n", msgEn)

	ctxZh := context.Background()
	ctxZh = i18n.SetLanguage(ctxZh, language.Chinese)

	msgZh := i18n.T(ctxZh, "validation_unknown_field")
	fmt.Printf("  Chinese: %s\n", msgZh)
	fmt.Println()

	// Example 2: Error translation with parameters
	fmt.Println("2. Error Translation with Parameters:")

	params := map[string]any{"min": 10}

	errorEn := i18n.TranslateError(ctxEn, "validation_min_text_constraint", "Must be at least {{.min}} character(s).", params)
	fmt.Printf("  English: %s\n", errorEn)

	errorZh := i18n.TranslateError(ctxZh, "validation_min_text_constraint", "必须至少包含{{.min}}个字符。", params)
	fmt.Printf("  Chinese: %s\n", errorZh)
	fmt.Println()

	// Example 3: Language detection from HTTP header
	fmt.Println("3. Language Detection from HTTP Header:")

	// Simulate HTTP request with Accept-Language header
	testCases := []struct {
		header string
		desc   string
	}{
		{"zh-CN,zh;q=0.9", "Chinese client"},
		{"en-US,en;q=0.8", "English client"},
		{"fr-FR,fr;q=0.7", "French client (fallback to English)"},
		{"", "No header (fallback to English)"},
	}

	for _, tc := range testCases {
		lang := i18n.DetectLanguage(tc.header)
		langName := lang.String()
		fmt.Printf("  %s: %s -> %s\n", tc.desc, tc.header, langName)
	}
	fmt.Println()

	// Example 4: HTTP middleware simulation
	fmt.Println("4. HTTP Middleware Simulation:")

	// Create a test request with Accept-Language header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// Detect language from request
	detectedLang := i18n.DetectLanguage(req.Header.Get("Accept-Language"))

	// Set language in context
	ctx := i18n.SetLanguage(req.Context(), detectedLang)
	req = req.WithContext(ctx)

	// Get language from context
	retrievedLang := i18n.GetLanguage(req.Context())
	fmt.Printf("  Request Accept-Language: %s\n", req.Header.Get("Accept-Language"))
	fmt.Printf("  Detected language: %s\n", detectedLang.String())
	fmt.Printf("  Retrieved from context: %s\n", retrievedLang.String())

	// Use the language for translation
	translatedMsg := i18n.T(req.Context(), "error_not_found")
	fmt.Printf("  Translated 'error_not_found': %s\n", translatedMsg)
	fmt.Println()

	// Example 5: API error responses
	fmt.Println("5. API Error Responses:")

	// Simulate API error creation with i18n
	apiErrorEn := i18n.TranslateError(ctxEn, "validation_invalid_token", "Invalid or expired token.", nil)
	apiErrorZh := i18n.TranslateError(ctxZh, "validation_invalid_token", "令牌无效或已过期。", nil)

	fmt.Printf("  English API error: %s\n", apiErrorEn)
	fmt.Printf("  Chinese API error: %s\n", apiErrorZh)
	fmt.Println()

	// Example 6: Common validation errors
	fmt.Println("6. Common Validation Errors:")

	commonErrors := []string{
		"validation_required",
		"validation_invalid_format",
		"validation_not_unique",
		"validation_invalid_url",
		"validation_invalid_json",
	}

	fmt.Println("  English translations:")
	for _, errCode := range commonErrors {
		msg := i18n.T(ctxEn, errCode)
		fmt.Printf("    %s: %s\n", errCode, msg)
	}

	fmt.Println("\n  Chinese translations:")
	for _, errCode := range commonErrors {
		msg := i18n.T(ctxZh, errCode)
		fmt.Printf("    %s: %s\n", errCode, msg)
	}
	fmt.Println()

	fmt.Println("=== Migration Guide ===")
	fmt.Println()
	fmt.Println("To migrate existing PocketBase code to use i18n:")
	fmt.Println()
	fmt.Println("1. Add language detection middleware to your router:")
	fmt.Println("   router.Use(apis.LanguageDetection())")
	fmt.Println()
	fmt.Println("2. Replace error creation methods:")
	fmt.Println("   - e.BadRequestError() -> e.BadRequestErrorWithContext()")
	fmt.Println("   - e.NotFoundError() -> e.NotFoundErrorWithContext()")
	fmt.Println("   - e.ForbiddenError() -> e.ForbiddenErrorWithContext()")
	fmt.Println("   - e.UnauthorizedError() -> e.UnauthorizedErrorWithContext()")
	fmt.Println("   - e.InternalServerError() -> e.InternalServerErrorWithContext()")
	fmt.Println()
	fmt.Println("3. For validation errors, ensure error codes are used as translation keys")
	fmt.Println("   Example: validation.NewError('validation_required', 'This field is required.')")
	fmt.Println()
	fmt.Println("4. Add more translations by calling message.SetString() in i18n.Init()")
	fmt.Println()
	fmt.Println("5. To add support for more languages:")
	fmt.Println("   - Add language tag to supportedLanguages in i18n.go")
	fmt.Println("   - Add translations using message.SetString()")
	fmt.Println("   - Update language detection if needed")
}
