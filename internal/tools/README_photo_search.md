# Photo Search Tool

这是一个基于 Pexels API 的图片搜索工具，可以根据关键词搜索高质量的图片并返回中等尺寸的图片链接。

## 功能特性

- 支持关键词搜索图片
- 返回中等尺寸 (medium) 的图片链接
- 包含图片的详细信息（摄影师、尺寸、描述等）
- 支持自定义返回图片数量（1-80张）
- 提供人类可读和JSON格式的双重输出

## 配置要求

在 `manifest/config/config.yaml` 中配置 Pexels API Key：

```yaml
ai:
  pexelsApiKey: "YOUR_PEXELS_API_KEY"
```

## API 参数

- `query` (必需): 搜索关键词，如 "nature", "city", "people"
- `per_page` (可选): 返回图片数量，默认5张，最多80张

## 使用示例

### 通过 AI 助手使用

用户可以直接向 AI 助手发送请求：

```
"帮我搜索一些自然风景的图片"
"找几张城市夜景的照片"
"搜索5张关于咖啡的图片"
```

### 直接调用示例

```json
{
  "query": "nature",
  "per_page": 3
}
```

## 返回格式

工具会返回包含以下信息的结果：

1. **人类可读格式**：
   - 搜索结果概述
   - 每张图片的描述、链接、摄影师和尺寸信息

2. **JSON 格式**：
   ```json
   [
     {
       "id": 3573351,
       "alt": "Brown Rocks During Golden Hour",
       "medium_url": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&h=350",
       "photographer": "Lukas Rodriguez",
       "width": 3066,
       "height": 3968
     }
   ]
   ```

## 错误处理

- 自动验证 API Key 配置
- 处理网络请求错误
- 验证搜索参数
- 处理 API 限制和配额

## 注意事项

1. 需要有效的 Pexels API Key
2. 遵守 Pexels API 使用条款
3. 图片链接有时效性，建议及时使用
4. API 有请求频率限制，请合理使用