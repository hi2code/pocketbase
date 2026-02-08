# PocketBase i18n 国际化实现指南

## 概述

本指南详细介绍了如何在PocketBase项目中实现国际化(i18n)支持，将错误消息从英文改为中文，并支持多语言。

## 实现内容

### 1. i18n包 (`tools/i18n/i18n.go`)
- 使用Go标准库的`golang.org/x/text/message`包
- 支持英语和中文两种语言
- 提供语言检测、翻译和上下文管理功能

### 2. 错误处理增强 (`tools/router/error.go`)
- 添加了支持上下文的错误创建函数：
  - `NewApiErrorWithContext()`
  - `NewNotFoundErrorWithContext()`
  - `NewBadRequestErrorWithContext()`
  - `NewForbiddenErrorWithContext()`
  - `NewUnauthorizedErrorWithContext()`
  - `NewInternalServerErrorWithContext()`
  - `NewTooManyRequestsErrorWithContext()`
  - `ToApiErrorWithContext()`

### 3. Event接口增强 (`tools/router/event.go`)
- 添加了支持i18n的错误方法：
  - `ErrorWithContext()`
  - `BadRequestErrorWithContext()`
  - `NotFoundErrorWithContext()`
  - `ForbiddenErrorWithContext()`
  - `UnauthorizedErrorWithContext()`
  - `TooManyRequestsErrorWithContext()`
  - `InternalServerErrorWithContext()`

### 4. 语言检测中间件 (`apis/middlewares.go`)
- 添加了`LanguageDetection()`中间件
- 自动从`Accept-Language`HTTP头检测用户语言
- 将检测到的语言设置到请求上下文中

### 5. 应用初始化 (`core/base.go`)
- 在`Bootstrap()`方法中初始化i18n系统
- 确保应用启动时加载所有翻译

## 翻译内容

### 错误消息翻译
已翻译以下类型的错误消息：

1. **通用字段错误** (100+条)
   - `validation_unknown_field`: "未知或无效的字段。"
   - `validation_invalid_value`: "无效的值。"
   - `validation_required`: "此字段为必填项。"

2. **HTTP API错误**
   - `error_not_found`: "请求的资源未找到。"
   - `error_bad_request`: "处理请求时出错。"
   - `error_forbidden`: "您无权执行此请求。"
   - `error_unauthorized`: "缺少或无效的身份验证。"

3. **验证错误**
   - 文本约束、格式验证、唯一性检查等
   - 文件上传、地理位置、JSON格式等

## 使用方法

### 1. 启用语言检测中间件
```go
router.Use(apis.LanguageDetection())
```

### 2. 使用支持i18n的错误方法
替换现有的错误创建方法：

```go
// 之前
return e.BadRequestError("Something went wrong", nil)

// 之后
return e.BadRequestErrorWithContext("Something went wrong", nil)
```

### 3. 验证错误自动翻译
验证错误会自动使用错误代码作为翻译键：
```go
// 错误代码会自动用作翻译键
validation.NewError("validation_required", "This field is required.")
```

## 添加新语言

### 1. 在i18n.go中添加语言支持
```go
var supportedLanguages = []language.Tag{
    language.English,
    language.Chinese,
    language.SimplifiedChinese,
    // 添加新语言
    language.Japanese,
    language.Korean,
}
```

### 2. 添加翻译
在`Init()`函数中添加翻译：
```go
message.SetString(language.Japanese, "validation_required", "このフィールドは必須です。")
message.SetString(language.Korean, "validation_required", "이 필드는 필수입니다.")
```

### 3. 更新默认语言（可选）
```go
var defaultLang = language.Japanese // 或根据需要更改
```

## 测试

### 运行i18n测试
```bash
go test ./tools/i18n/... -v
```

### 查看示例
```bash
go run examples/i18n_example.go
```

## 向后兼容性

### 保持向后兼容
- 原有的错误创建方法仍然可用
- 新的`*WithContext()`方法提供i18n支持
- 应用可以逐步迁移到i18n版本

### 迁移路径
1. 首先添加语言检测中间件
2. 逐步将错误创建方法替换为`*WithContext()`版本
3. 确保所有验证错误使用正确的错误代码

## 性能考虑

### 轻量级实现
- 使用Go标准库，无额外依赖
- 翻译在内存中，访问速度快
- 语言检测使用高效的匹配算法

### 内存使用
- 翻译文本存储在内存中
- 支持的语言越多，内存使用越大
- 建议只添加实际需要的语言

## 扩展性

### 添加更多错误类型
1. 在`i18n.go`的`Init()`函数中添加新翻译
2. 确保错误代码格式一致：`validation_*`或`error_*`

### 支持动态语言切换
- 可以通过URL参数、cookie或用户设置覆盖语言检测
- 实现自定义中间件来支持这些功能

## 注意事项

### 错误代码作为翻译键
- 所有验证错误应使用错误代码作为翻译键
- 错误代码格式：`validation_*`（验证错误）或`error_*`（HTTP错误）

### 参数化消息
- 支持参数化消息：`Must be at least {{.min}} character(s).`
- 在翻译时保持参数占位符

### 默认语言回退
- 如果翻译未找到，使用默认英文消息
- 确保所有错误都有英文翻译作为回退

## 总结

本实现为PocketBase提供了完整的i18n支持，具有以下特点：

1. **完整覆盖**：翻译了所有核心错误消息
2. **易于使用**：简单的API，向后兼容
3. **高性能**：基于标准库，内存效率高
4. **可扩展**：易于添加新语言和翻译
5. **自动化**：语言检测和错误翻译自动处理

通过此实现，PocketBase可以更好地服务全球用户，提供本地化的错误消息和更好的用户体验。