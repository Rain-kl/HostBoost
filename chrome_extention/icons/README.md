# 图标说明

## SVG 图标已创建
当前项目包含以下 SVG 格式的图标:
- icon16.svg (16x16)
- icon48.svg (48x48)  
- icon128.svg (128x128)

## 转换为 PNG (可选)

虽然 Chrome 扩展支持 SVG 图标,但如果你需要 PNG 格式,可以使用以下方法:

### 方法1: 使用在线工具
访问 https://cloudconvert.com/svg-to-png 或 https://svgtopng.com/ 
上传 SVG 文件并转换为对应尺寸的 PNG

### 方法2: 使用命令行工具 (需要安装 ImageMagick 或 Inkscape)
```bash
# 使用 ImageMagick
convert icon16.svg -resize 16x16 icon16.png
convert icon48.svg -resize 48x48 icon48.png
convert icon128.svg -resize 128x128 icon128.png

# 或使用 Inkscape
inkscape icon16.svg -w 16 -h 16 -o icon16.png
inkscape icon48.svg -w 48 -h 48 -o icon48.png
inkscape icon128.svg -w 128 -h 128 -o icon128.png
```

## 设计说明
图标设计采用了渐变紫色背景配合白色盾牌图案,象征安全和保护。
中心的绿色圆点表示连接状态,整体风格现代简约。
