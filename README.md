# ExchangeApp 汇率兑换与资讯示例项目

一个全栈示例项目：**Go + Gin + Gorm + Redis + MySQL** 后端，**Vue 3 + Vite + TypeScript + Pinia + Element Plus** 前端。功能包括用户注册 / 登录（JWT 鉴权）、汇率管理、文章发布与点赞（Redis 缓存）。

> 本项目源自 B 站 [InkkaPlum 频道](https://space.bilibili.com/290859233) 的 Gin / Gorm / Vue / Redis / MySQL 教程，已与原作者仓库解绑，仅保留学习用途的代码骨架。

## 技术栈

### 后端 `Exchangeapp_backend`
- Go 1.22
- Gin（HTTP 框架）
- Gorm（ORM / MySQL）
- go-redis（缓存）
- JWT（鉴权）

### 前端 `Exchangeapp_frontend`
- Vue 3 + Vite
- TypeScript
- Pinia（状态管理）
- Element Plus（UI）
- axios（HTTP）

## 目录结构

```
.
├── Exchangeapp_backend/   后端
│   ├── config/            配置、数据库 / Redis 初始化
│   ├── controllers/       控制器
│   ├── global/            全局变量
│   ├── middlewares/       中间件（JWT 鉴权）
│   ├── models/            数据模型
│   ├── router/            路由
│   └── utils/             工具（JWT、密码哈希）
└── Exchangeapp_frontend/  前端
    └── src/
        ├── views/         页面
        ├── components/    组件
        ├── store/         Pinia store
        └── router/        前端路由
```

## 运行

### 后端
1. 准备 MySQL 与 Redis，在 `Exchangeapp_backend/config/config.yml` 中配置连接信息。
2. 启动：

   ```bash
   cd Exchangeapp_backend
   go run .
   ```

   默认监听 `:3000`。

### 前端

```bash
cd Exchangeapp_frontend
npm install
npm run dev
```

前端请求会指向后端 `http://localhost:3000/api`（见 `src/axios.ts`）。

## 注意

- 后端与前端中的密钥、数据库口令为示例占位，正式使用请改用环境变量或密钥管理服务。
