package i18n

import (
	"context"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Default language
var defaultLang = language.English

// Supported languages
var supportedLanguages = []language.Tag{
	language.English,           // en
	language.Chinese,           // zh
	language.SimplifiedChinese, // zh-Hans
}

// Matcher for language negotiation
var matcher = language.NewMatcher(supportedLanguages)

// Context key for language
type langKey struct{}

// SetLanguage sets the language in context
func SetLanguage(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, langKey{}, lang)
}

// GetLanguage gets the language from context, returns default if not set
func GetLanguage(ctx context.Context) language.Tag {
	if lang, ok := ctx.Value(langKey{}).(language.Tag); ok {
		return lang
	}
	return defaultLang
}

// DetectLanguage detects language from Accept-Language header
func DetectLanguage(acceptLang string) language.Tag {
	tags, _, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil || len(tags) == 0 {
		return defaultLang
	}

	matched, _, _ := matcher.Match(tags...)
	return matched
}

// T translates a message using the language from context
func T(ctx context.Context, key string, args ...any) string {
	lang := GetLanguage(ctx)
	p := message.NewPrinter(lang)

	// Try to get translation, fallback to key if not found
	return p.Sprintf(key, args...)
}

// TWithLang translates a message using specified language
func TWithLang(lang language.Tag, key string, args ...any) string {
	p := message.NewPrinter(lang)
	return p.Sprintf(key, args...)
}

// Init initializes i18n with translations
func Init() error {
	// English translations (default)
	message.SetString(language.English, "validation_unknown_field", "Unknown or invalid field.")
	message.SetString(language.English, "validation_invalid_field_value", "Invalid field value.")
	message.SetString(language.English, "validation_must_be_system_and_hidden", `The field must be marked as "System" and "Hidden".`)
	message.SetString(language.English, "validation_must_be_system", `The field must be marked as "System".`)
	message.SetString(language.English, "validation_invalid_value", "Invalid value.")
	message.SetString(language.English, "validation_required", "This field is required.")
	message.SetString(language.English, "validation_invalid_format", "Invalid value format.")
	message.SetString(language.English, "validation_invalid_json", "Must be a valid json value.")
	message.SetString(language.English, "validation_invalid_url", "Must be a valid url.")
	message.SetString(language.English, "validation_invalid_token", "Invalid or expired token.")
	message.SetString(language.English, "validation_not_unique", "Value must be unique.")
	message.SetString(language.English, "validation_values_mismatch", "Values don't match.")
	message.SetString(language.English, "validation_min_text_constraint", "Must be at least %d character(s).")
	message.SetString(language.English, "validation_max_text_constraint", "Must be no more than %d character(s).")
	message.SetString(language.English, "validation_pk_change", "The record primary key cannot be changed.")
	message.SetString(language.English, "validation_pk_invalid", "The record primary key is invalid or already exists.")
	message.SetString(language.English, "validation_forbidden_pk_character", "'%s' is not a valid primary key character.")
	message.SetString(language.English, "validation_reserved_pk", "The primary key '%s' is reserved and cannot be used.")
	message.SetString(language.English, "validation_not_enough_values", "Select at least %d.")
	message.SetString(language.English, "validation_too_many_values", "Select no more than %d.")
	message.SetString(language.English, "validation_missing_rel_collection", "Relation connection is missing or cannot be accessed.")
	message.SetString(language.English, "validation_missing_rel_records", "Failed to find all relation records with the provided ids.")
	message.SetString(language.English, "validation_invalid_file", "Invalid new files: %s.")
	message.SetString(language.English, "validation_too_many_files", "The maximum allowed files is %d.")
	message.SetString(language.English, "validation_invalid_latitude", "Latitude must be between -90 and 90 degrees.")
	message.SetString(language.English, "validation_invalid_longitude", "Longitude must be between -180 and 180 degrees.")
	message.SetString(language.English, "validation_invalid_cron", "%s")
	message.SetString(language.English, "validation_conflicting_rate_limit_rule", "Rate limit rule configuration with label %s already exists or conflicts with another rule.")
	message.SetString(language.English, "validation_invalid_auth_collection", "Must be a valid auth collection id or name.")
	message.SetString(language.English, "validation_invalid_autogenerate_pattern", "%s")
	message.SetString(language.English, "validation_collection_name_exists", "Collection name must be unique (case insensitive).")
	message.SetString(language.English, "validation_collection_name_id_duplicate", "The name must not match an existing collection id.")
	message.SetString(language.English, "validation_collection_name_invalid", "The name shouldn't match with an existing internal table.")
	message.SetString(language.English, "validation_collection_system_name_change", "System collection name cannot be changed.")
	message.SetString(language.English, "validation_collection_system_flag_change", "System collection state cannot be changed.")
	message.SetString(language.English, "validation_collection_type_change", "Collection type cannot be changed.")
	message.SetString(language.English, "validation_missing_primary_key", `Missing or invalid "id" PK field.`)
	message.SetString(language.English, "validation_missing_password_field", `System "password" field is required.`)
	message.SetString(language.English, "validation_missing_tokenKey_field", `System "tokenKey" field is required.`)
	message.SetString(language.English, "validation_missing_email_field", `System "email" field is required.`)
	message.SetString(language.English, "validation_missing_emailVisibility_field", `System "emailVisibility" field is required.`)
	message.SetString(language.English, "validation_missing_verified_field", `System "verified" field is required.`)
	message.SetString(language.English, "validation_system_field_change", "System fields cannot be deleted or renamed.")
	message.SetString(language.English, "validation_missing_field", "Invalid or missing field %s.")
	message.SetString(language.English, "validation_missing_unique_constraint", "The field %s doesn't have a UNIQUE constraint.")
	message.SetString(language.English, "validation_invalid_rule", "Invalid rule. Raw error: %s")
	message.SetString(language.English, "validation_collection_system_rule_change", "System collection API rule cannot be changed.")
	message.SetString(language.English, "validation_invalid_old_password", "Missing or invalid old password.")
	message.SetString(language.English, "validation_invalid_auth_id", "Invalid or duplicated auth record id.")
	message.SetString(language.English, "validation_invalid_provider", "Provider with name %s is missing or is not enabled.")
	message.SetString(language.English, "validation_invalid_token_claims", "Missing email token claim.")
	message.SetString(language.English, "validation_token_collection_mismatch", "The provided token is for different auth collection.")
	message.SetString(language.English, "validation_token_email_mismatch", "The record email doesn't match with the requested token claims.")
	message.SetString(language.English, "validation_invalid_new_email", "Invalid new email address.")
	message.SetString(language.English, "validation_invalid_password", "Missing or invalid auth record password.")
	message.SetString(language.English, "validation_invalid_token_payload", "Invalid token payload - newEmail must be set.")
	message.SetString(language.English, "validation_existing_token_email", "The new email address is already registered: %s")
	message.SetString(language.English, "validation_invalid_collection_id", "Missing or invalid collection.")
	message.SetString(language.English, "validation_invalid_collection", "Missing or invalid collection.")
	message.SetString(language.English, "validation_invalid_record", "Missing or invalid record.")
	message.SetString(language.English, "validation_invalid_or_existing_id", "The model id is invalid or already exists.")
	message.SetString(language.English, "validation_invalid_regex", "%s")
	message.SetString(language.English, "validation_unsupported_value_type", "Invalid or unsupported value type.")
	message.SetString(language.English, "validation_unsupported_composite_pk", "Composite PKs are not supported and the collection must have only 1 PK.")

	// HTTP error messages
	message.SetString(language.English, "error_not_found", "The requested resource wasn't found.")
	message.SetString(language.English, "error_bad_request", "Something went wrong while processing your request.")
	message.SetString(language.English, "error_forbidden", "You are not allowed to perform this request.")
	message.SetString(language.English, "error_unauthorized", "Missing or invalid authentication.")
	message.SetString(language.English, "error_internal_server", "Something went wrong while processing your request.")
	message.SetString(language.English, "error_too_many_requests", "Too Many Requests.")

	// Chinese translations
	message.SetString(language.Chinese, "validation_unknown_field", "未知或无效的字段。")
	message.SetString(language.Chinese, "validation_invalid_field_value", "字段值无效。")
	message.SetString(language.Chinese, "validation_must_be_system_and_hidden", `字段必须标记为"系统"和"隐藏"。`)
	message.SetString(language.Chinese, "validation_must_be_system", `字段必须标记为"系统"。`)
	message.SetString(language.Chinese, "validation_invalid_value", "无效的值。")
	message.SetString(language.Chinese, "validation_required", "此字段为必填项。")
	message.SetString(language.Chinese, "validation_invalid_format", "值格式无效。")
	message.SetString(language.Chinese, "validation_invalid_json", "必须是有效的JSON值。")
	message.SetString(language.Chinese, "validation_invalid_url", "必须是有效的URL。")
	message.SetString(language.Chinese, "validation_invalid_token", "令牌无效或已过期。")
	message.SetString(language.Chinese, "validation_not_unique", "值必须是唯一的。")
	message.SetString(language.Chinese, "validation_values_mismatch", "值不匹配。")
	message.SetString(language.Chinese, "validation_min_text_constraint", "必须至少包含%d个字符。")
	message.SetString(language.Chinese, "validation_max_text_constraint", "不能超过%d个字符。")
	message.SetString(language.Chinese, "validation_pk_change", "记录主键不能更改。")
	message.SetString(language.Chinese, "validation_pk_invalid", "记录主键无效或已存在。")
	message.SetString(language.Chinese, "validation_forbidden_pk_character", "'%s'不是有效的主键字符。")
	message.SetString(language.Chinese, "validation_reserved_pk", "主键'%s'是保留的，不能使用。")
	message.SetString(language.Chinese, "validation_not_enough_values", "请至少选择%d个。")
	message.SetString(language.Chinese, "validation_too_many_values", "请选择不超过%d个。")
	message.SetString(language.Chinese, "validation_missing_rel_collection", "关系连接缺失或无法访问。")
	message.SetString(language.Chinese, "validation_missing_rel_records", "未能找到提供的ID对应的所有关系记录。")
	message.SetString(language.Chinese, "validation_invalid_file", "无效的新文件：%s。")
	message.SetString(language.Chinese, "validation_too_many_files", "允许的最大文件数为%d。")
	message.SetString(language.Chinese, "validation_invalid_latitude", "纬度必须在-90到90度之间。")
	message.SetString(language.Chinese, "validation_invalid_longitude", "经度必须在-180到180度之间。")
	message.SetString(language.Chinese, "validation_invalid_cron", "%s")
	message.SetString(language.Chinese, "validation_conflicting_rate_limit_rule", "标签为%s的速率限制规则配置已存在或与其他规则冲突。")
	message.SetString(language.Chinese, "validation_invalid_auth_collection", "必须是有效的认证集合ID或名称。")
	message.SetString(language.Chinese, "validation_invalid_autogenerate_pattern", "%s")
	message.SetString(language.Chinese, "validation_collection_name_exists", "集合名称必须是唯一的（不区分大小写）。")
	message.SetString(language.Chinese, "validation_collection_name_id_duplicate", "名称不能与现有集合ID匹配。")
	message.SetString(language.Chinese, "validation_collection_name_invalid", "名称不应与现有内部表匹配。")
	message.SetString(language.Chinese, "validation_collection_system_name_change", "系统集合名称不能更改。")
	message.SetString(language.Chinese, "validation_collection_system_flag_change", "系统集合状态不能更改。")
	message.SetString(language.Chinese, "validation_collection_type_change", "集合类型不能更改。")
	message.SetString(language.Chinese, "validation_missing_primary_key", `缺少或无效的"id"主键字段。`)
	message.SetString(language.Chinese, "validation_missing_password_field", `系统"password"字段是必需的。`)
	message.SetString(language.Chinese, "validation_missing_tokenKey_field", `系统"tokenKey"字段是必需的。`)
	message.SetString(language.Chinese, "validation_missing_email_field", `系统"email"字段是必需的。`)
	message.SetString(language.Chinese, "validation_missing_emailVisibility_field", `系统"emailVisibility"字段是必需的。`)
	message.SetString(language.Chinese, "validation_missing_verified_field", `系统"verified"字段是必需的。`)
	message.SetString(language.Chinese, "validation_system_field_change", "系统字段不能被删除或重命名。")
	message.SetString(language.Chinese, "validation_missing_field", "无效或缺失的字段%s。")
	message.SetString(language.Chinese, "validation_missing_unique_constraint", "字段%s没有UNIQUE约束。")
	message.SetString(language.Chinese, "validation_invalid_rule", "无效的规则。原始错误：%s")
	message.SetString(language.Chinese, "validation_collection_system_rule_change", "系统集合API规则不能更改。")
	message.SetString(language.Chinese, "validation_invalid_old_password", "缺少或无效的旧密码。")
	message.SetString(language.Chinese, "validation_invalid_auth_id", "无效或重复的认证记录ID。")
	message.SetString(language.Chinese, "validation_invalid_provider", "名称为%s的提供程序缺失或未启用。")
	message.SetString(language.Chinese, "validation_invalid_token_claims", "缺少电子邮件令牌声明。")
	message.SetString(language.Chinese, "validation_token_collection_mismatch", "提供的令牌用于不同的认证集合。")
	message.SetString(language.Chinese, "validation_token_email_mismatch", "记录电子邮件与请求的令牌声明不匹配。")
	message.SetString(language.Chinese, "validation_invalid_new_email", "无效的新电子邮件地址。")
	message.SetString(language.Chinese, "validation_invalid_password", "缺少或无效的认证记录密码。")
	message.SetString(language.Chinese, "validation_invalid_token_payload", "无效的令牌负载 - 必须设置newEmail。")
	message.SetString(language.Chinese, "validation_existing_token_email", "新电子邮件地址已注册：%s")
	message.SetString(language.Chinese, "validation_invalid_collection_id", "缺少或无效的集合。")
	message.SetString(language.Chinese, "validation_invalid_collection", "缺少或无效的集合。")
	message.SetString(language.Chinese, "validation_invalid_record", "缺少或无效的记录。")
	message.SetString(language.Chinese, "validation_invalid_or_existing_id", "模型ID无效或已存在。")
	message.SetString(language.Chinese, "validation_invalid_regex", "%s")
	message.SetString(language.Chinese, "validation_unsupported_value_type", "无效或不支持的值类型。")
	message.SetString(language.Chinese, "validation_unsupported_composite_pk", "不支持复合主键，集合必须只有1个主键。")

	// HTTP error messages in Chinese
	message.SetString(language.Chinese, "error_not_found", "请求的资源未找到。")
	message.SetString(language.Chinese, "error_bad_request", "处理请求时出错。")
	message.SetString(language.Chinese, "error_forbidden", "您无权执行此请求。")
	message.SetString(language.Chinese, "error_unauthorized", "缺少或无效的身份验证。")
	message.SetString(language.Chinese, "error_internal_server", "处理请求时出错。")
	message.SetString(language.Chinese, "error_too_many_requests", "请求过多。")

	return nil
}

// TranslateError translates an error message using context
func TranslateError(ctx context.Context, code string, defaultMsg string, params map[string]any) string {
	lang := GetLanguage(ctx)

	// Convert params to interface slice for formatting
	var args []any
	if params != nil {
		// For messages with named parameters, we need to handle them specially
		// For now, just use the default message formatting
		args = make([]any, 0, len(params))
		for _, v := range params {
			args = append(args, v)
		}
	}

	translated := TWithLang(lang, code, args...)

	// If translation returns the key (no translation found), use default message
	if translated == code {
		// Format default message with params if available
		if len(args) > 0 {
			return fmt.Sprintf(defaultMsg, args...)
		}
		return defaultMsg
	}

	return translated
}

// Helper function to extract error code from error
func ExtractErrorCode(err error) string {
	if se, ok := err.(interface{ Code() string }); ok {
		return se.Code()
	}
	return ""
}

// Helper function to extract error message from error
func ExtractErrorMessage(err error) string {
	if se, ok := err.(error); ok {
		return se.Error()
	}
	return ""
}
