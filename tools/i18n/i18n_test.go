package i18n

import (
	"context"
	"testing"

	"golang.org/x/text/language"
)

func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name       string
		acceptLang string
		want       language.Tag
	}{
		{
			name:       "Chinese",
			acceptLang: "zh-CN,zh;q=0.9,en;q=0.8",
			want:       language.Chinese,
		},
		{
			name:       "English",
			acceptLang: "en-US,en;q=0.9",
			want:       language.English,
		},
		{
			name:       "Simplified Chinese",
			acceptLang: "zh-Hans-CN,zh-Hans;q=0.8",
			want:       language.SimplifiedChinese,
		},
		{
			name:       "Empty header",
			acceptLang: "",
			want:       language.English,
		},
		{
			name:       "Invalid header",
			acceptLang: "invalid",
			want:       language.English,
		},
	}

	Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectLanguage(tt.acceptLang)
			// Compare base language (ignore regional variants)
			gotBase, _ := got.Base()
			wantBase, _ := tt.want.Base()
			if gotBase != wantBase {
				t.Errorf("DetectLanguage(%q) = %v (base: %v), want %v (base: %v)", tt.acceptLang, got, gotBase, tt.want, wantBase)
			}
		})
	}
}

func TestTranslate(t *testing.T) {
	Init()

	ctx := context.Background()

	// Test English translation
	ctxEn := SetLanguage(ctx, language.English)
	gotEn := T(ctxEn, "validation_unknown_field")
	wantEn := "Unknown or invalid field."
	if gotEn != wantEn {
		t.Errorf("T(English, 'validation_unknown_field') = %q, want %q", gotEn, wantEn)
	}

	// Test Chinese translation
	ctxZh := SetLanguage(ctx, language.Chinese)
	gotZh := T(ctxZh, "validation_unknown_field")
	wantZh := "未知或无效的字段。"
	if gotZh != wantZh {
		t.Errorf("T(Chinese, 'validation_unknown_field') = %q, want %q", gotZh, wantZh)
	}

	// Test fallback to key when translation not found
	gotFallback := T(ctxEn, "non_existent_key")
	if gotFallback != "non_existent_key" {
		t.Errorf("T(English, 'non_existent_key') = %q, want 'non_existent_key'", gotFallback)
	}
}

func TestTranslateError(t *testing.T) {
	Init()

	ctx := context.Background()

	// Test English error translation
	ctxEn := SetLanguage(ctx, language.English)
	gotEn := TranslateError(ctxEn, "validation_invalid_value", "Invalid value.", nil)
	wantEn := "Invalid value."
	if gotEn != wantEn {
		t.Errorf("TranslateError(English, 'validation_invalid_value') = %q, want %q", gotEn, wantEn)
	}

	// Test Chinese error translation
	ctxZh := SetLanguage(ctx, language.Chinese)
	gotZh := TranslateError(ctxZh, "validation_invalid_value", "Invalid value.", nil)
	wantZh := "无效的值。"
	if gotZh != wantZh {
		t.Errorf("TranslateError(Chinese, 'validation_invalid_value') = %q, want %q", gotZh, wantZh)
	}

	// Test with parameters
	params := map[string]any{"min": 5}
	gotEnWithParams := TranslateError(ctxEn, "validation_min_text_constraint", "Must be at least {{.min}} character(s).", params)
	wantEnWithParams := "Must be at least 5 character(s)."
	if gotEnWithParams != wantEnWithParams {
		t.Errorf("TranslateError(English, 'validation_min_text_constraint') = %q, want %q", gotEnWithParams, wantEnWithParams)
	}

	gotZhWithParams := TranslateError(ctxZh, "validation_min_text_constraint", "必须至少包含{{.min}}个字符。", params)
	wantZhWithParams := "必须至少包含5个字符。"
	if gotZhWithParams != wantZhWithParams {
		t.Errorf("TranslateError(Chinese, 'validation_min_text_constraint') = %q, want %q", gotZhWithParams, wantZhWithParams)
	}
}

func TestGetLanguage(t *testing.T) {
	ctx := context.Background()

	// Test default language
	got := GetLanguage(ctx)
	if got != language.English {
		t.Errorf("GetLanguage() with empty context = %v, want %v", got, language.English)
	}

	// Test with language set
	ctxWithLang := SetLanguage(ctx, language.Chinese)
	got = GetLanguage(ctxWithLang)
	if got != language.Chinese {
		t.Errorf("GetLanguage() with Chinese context = %v, want %v", got, language.Chinese)
	}
}
