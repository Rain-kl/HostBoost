# HostBoost VPN - Chrome 扩展

一个基于 Manifest V3 的 Chrome 扩展程序,仿照 AdGuard VPN 的界面风格,用于一键启用代理服务。

## 功能特性

- 🎨 **现代化界面**: 仿照 AdGuard VPN 的设计风格,渐变紫色主题,简洁美观
- 🔄 **一键切换**: 通过开关按钮快速启用/禁用代理服务
- 🌐 **域名识别**: 自动获取当前标签页的域名信息
- 💾 **状态保持**: 使用 Chrome Storage API 保存启用状态
- ⚡ **实时通信**: 向本地服务器发送 POST 请求控制代理

## 项目结构

```
chrome_extention/
├── manifest.json          # Manifest V3 配置文件
├── popup.html            # 弹出窗口页面
├── popup.css             # 样式文件
├── popup.js              # 逻辑脚本
├── icons/                # 图标资源
│   ├── icon16.svg
│   ├── icon48.svg
│   ├── icon128.svg
│   └── README.md
└── README.md            # 说明文档
```

## 安装方法

### 1. 开发者模式加载

1. 打开 Chrome 浏览器,访问 `chrome://extensions/`
2. 在右上角启用 **"开发者模式"**
3. 点击 **"加载已解压的扩展程序"**
4. 选择本项目的 `chrome_extention` 文件夹
5. 扩展将出现在工具栏中

### 2. 后端服务要求

扩展需要本地运行一个后端服务,监听 `127.0.0.1:8003`,并提供以下接口:

#### POST /enable
启用代理服务

**请求体:**
```json
{
  "domain": "example.com",
  "timestamp": "2025-11-01T12:00:00.000Z"
}
```

#### POST /disable (可选)
禁用代理服务

**请求体:**
```json
{
  "domain": "example.com",
  "timestamp": "2025-11-01T12:00:00.000Z"
}
```

## 使用说明

1. **查看当前域名**: 打开扩展弹窗,会自动显示当前标签页的域名
2. **启用服务**: 点击切换开关,扩展会向 `http://127.0.0.1:8003/enable` 发送 POST 请求
3. **状态显示**: 
   - 🔴 已断开: 服务未启用
   - 🟢 已连接: 服务已启用
   - ⏳ 连接中: 正在发送请求
4. **禁用服务**: 再次点击开关即可禁用

## 技术栈

- **Manifest V3**: Chrome 扩展的最新规范
- **原生 JavaScript**: 无需额外依赖
- **Chrome APIs**: 
  - `chrome.tabs`: 获取标签页信息
  - `chrome.storage`: 保存状态数据
- **Fetch API**: 与后端服务通信

## 权限说明

扩展申请了以下权限:
- `activeTab`: 获取当前活动标签页的域名信息
- `storage`: 保存启用状态
- `tabs`: 监听标签页切换和更新
- `host_permissions`: 允许访问 `http://127.0.0.1:8003/*`

## 开发说明

### 修改服务器地址

如需修改后端服务地址,请编辑以下文件:

1. **manifest.json**: 修改 `host_permissions`
2. **popup.js**: 修改 `sendEnableRequest` 和 `sendDisableRequest` 函数中的 URL

### 自定义样式

编辑 `popup.css` 文件可以修改界面样式:
- 渐变背景色
- 开关按钮颜色
- 状态指示器动画

## 注意事项

⚠️ **重要提示:**
- 确保后端服务已启动并监听在 `127.0.0.1:8003`
- Chrome 可能会阻止对 localhost 的请求,需要在 manifest.json 中正确配置权限
- 扩展仅在有效的网页标签页中工作,无法在 `chrome://` 等特殊页面使用

## 后续优化建议

- [ ] 添加更多代理服务器选项
- [ ] 实现域名白名单/黑名单功能
- [ ] 添加连接统计和日志记录
- [ ] 支持自定义服务器地址配置
- [ ] 添加快捷键支持

## License

MIT License

## 作者

HostBoost Team

HostBoost 的 Chrome 扩展程序 (CRX) 用于增强浏览器与 HostBoost 服务的集成，提供更便捷的访问和管理功能。

## 功能

1. Boost 按钮: 点击则发送域名到 HostBoost 提供加速服务, 取消则移除加速服务。
