## Agent架构
<img src="./resource/public/resource/image/rctAgent.png" alt="架构图" width="50%">

### 说明
前后端分离项目
前端:
1. Vue3 项目
2. Axios 请求库
后端
1. GoFrame 框架
2. Eino 框架

# 运行流程

## 1. 下载项目
```bash
git clone https://github.com/ipvoov/Agent-go.git
```

## 2. 安装 milvus 向量知识库
```bash
docker-compose up -d
```

## 3. 修改配置文件

修改manifest/config/config.yaml中的apiKey 和你想使用的模型
想要使用搜索和地图工具,要修改 api

## 4. 启动项目

在项目根目录下执行：
```bash
go run main.go
```

如果你需要自动编译，可以执行：
```bash
gf run main.go
```

## 5.前端运行流程
```bash
cd /agent-frontend && npm i && npm run dev
```

## 2.前端
## 2.1 前端展示
<img src="./resource/public/resource/image/homepage_L.png" alt="主页图" width="50%">

## 2.2 Ai思考过程
### 2.2.1 对话框
<div style="display: flex; justify-content: space-around; flex-wrap: wrap;">
<img src="./resource/public/resource/image/answer_L.png" alt="对话图" style="width: 40%; margin: 10px;">
<img src="./resource/public/resource/image/think_L.png" alt="Ai思考过程展示" style="width: 40%; margin: 10px;">
</div>

### 2.2.2 Ai思考过程与调用tool过程与生成的PDF
<div style="display: flex; justify-content: space-around; flex-wrap: wrap;">
<img src="./resource/public/resource/image/thinking.png" alt="Ai思考过程" style="width: 40%; margin: 10px;">
<img src="./resource/public/resource/image/ai_PDF.png" alt="Ai_PDF" style="width: 40%; margin: 10px;">
</div>




