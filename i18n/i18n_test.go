package i18n_test

import (
	"testing"

	"github.com/BounkhongDev/bkgo/i18n"
)

func TestTranslate_KnownCode(t *testing.T) {
	tests := []struct {
		locale i18n.Locale
		code   string
		want   string
	}{
		{i18n.EN, "NOT_FOUND", "Resource not found"},
		{i18n.LO, "NOT_FOUND", "ບໍ່ພົບຂໍ້ມູນ"},
		{i18n.TH, "NOT_FOUND", "ไม่พบข้อมูล"},
		{i18n.ZH, "NOT_FOUND", "未找到资源"},
		{i18n.EN, "UNAUTHORIZED", "Unauthorized"},
		{i18n.LO, "UNAUTHORIZED", "ບໍ່ມີສິດເຂົ້າໃຊ້ງານ"},
	}

	for _, tc := range tests {
		t.Run(string(tc.locale)+"/"+tc.code, func(t *testing.T) {
			got := i18n.Translate(tc.locale, tc.code)
			if got != tc.want {
				t.Errorf("want %q got %q", tc.want, got)
			}
		})
	}
}

func TestTranslate_FallbackToEnglish(t *testing.T) {
	// "ja" (Japanese) is not registered — should fall back to English
	got := i18n.Translate("ja", "NOT_FOUND")
	if got != "Resource not found" {
		t.Errorf("expected English fallback, got %q", got)
	}
}

func TestTranslate_UnknownCode(t *testing.T) {
	got := i18n.Translate(i18n.EN, "TOTALLY_UNKNOWN_CODE")
	if got != "" {
		t.Errorf("expected empty string for unknown code, got %q", got)
	}
}

func TestRegister_CustomLanguage(t *testing.T) {
	i18n.Register("fr", map[string]string{
		"NOT_FOUND": "Ressource introuvable",
	})

	got := i18n.Translate("fr", "NOT_FOUND")
	if got != "Ressource introuvable" {
		t.Errorf("want 'Ressource introuvable' got %q", got)
	}
}

func TestRegister_OverrideMessage(t *testing.T) {
	i18n.Register(i18n.EN, map[string]string{
		"CUSTOM_CODE": "Custom message",
	})

	got := i18n.Translate(i18n.EN, "CUSTOM_CODE")
	if got != "Custom message" {
		t.Errorf("want 'Custom message' got %q", got)
	}
}

func TestFromHeader(t *testing.T) {
	tests := []struct {
		header string
		want   i18n.Locale
	}{
		{"lo", i18n.LO},
		{"lo-LA", i18n.LO},                    // full tag falls back to prefix
		{"th", i18n.TH},
		{"zh-CN", i18n.ZH},
		{"en-US", i18n.EN},
		{"", i18n.EN},                          // empty → EN
		{"xx-XX", i18n.EN},                     // unknown → EN
		{"lo-LA,en;q=0.9", i18n.LO},            // multi-value header
		{"th;q=0.8,en;q=0.5", i18n.TH},         // with quality weights
	}

	for _, tc := range tests {
		t.Run(tc.header, func(t *testing.T) {
			got := i18n.FromHeader(tc.header)
			if got != tc.want {
				t.Errorf("header %q: want %q got %q", tc.header, tc.want, got)
			}
		})
	}
}

func TestSupported(t *testing.T) {
	locales := i18n.Supported()
	if len(locales) < 4 {
		t.Errorf("expected at least 4 built-in locales, got %d", len(locales))
	}
}
