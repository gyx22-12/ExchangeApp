# ExchangeApp 汇率兑换与资讯示例项目

一个全栈示例项目：**Go + Gin + Gorm + Redis + MySQL** 后端，**Vue 3 + Vite + TypeScript + Pinia + Element Plus** 前端。功能包括用户注册 / 登录（JWT 鉴权）、汇率管理、文章发布与点赞（Redis 缓存）。

> 本项目源自 B 站 [InkkaPlum 频道](https://space.bilibili.com/290859233) 的 Gin / Gorm / Vue / Redis / MySQL 教程，已与原作者仓库解绑，进行了大量优化和升级。

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

## 优化部分
   一、功能正确性

  1. 点赞去重（0fbae8d）
  - 问题：原实现用计数器，同一用户可以无限点赞。
  - 改法：Redis Set（SADD/SCARD）记录每个文章赞过的用户名，SADD 返回 0 说明已点过 →
    返回 409「你已经点过赞了」。
  - 讲点：为什么用 Set 而不是 INCR 计数器——Set 天然去重、且能拿到「谁赞了」。

  2. 注册用户名查重（cb65b3d）
  - 问题：重复用户名会撞 Username 的 unique 约束，后端直接 500，前端拿不到有用提示。
  - 改法：注册时先 First 查同名用户，存在则 409「username already exists」；DB 上的
    unique 约束保留做并发兜底（先查后插有竞态窗口）。
  - 讲点：先查后插 vs 依赖 DB 约束的取舍，以及并发下为什么还要留约束。

  二、代码质量（消除无效逻辑）

  3 & 4. 移除 Find 上的 ErrRecordNotFound 死判断（c56255c 的 GetArticles、8924384 的
  GetExchangeRates）
  - 问题：Find(&slice) 查空表返回空切片 + nil error，永远不会返回
    gorm.ErrRecordNotFound——那个 if errors.Is(...) 是死代码。
  - 改法：Find 的 err 分支简化为直接 InternalServerError。
  - 讲点：GORM 语义——First/Take/Last 查不到才返回 ErrRecordNotFound，Find
    不返回。这是很能体现「读得懂库语义」的一个点。

  三、性能 / 架构

  5. GetArticleByID 补 Redis 缓存（c56255c）
  - 问题：列表有 cache-aside，但详情接口每次直查 DB。
  - 改法：article:{id} 缓存，10 分钟 TTL，命中反序列化返回、未命中查库回填，和
    GetArticles 保持一致。

  6. AutoMigrate 从请求处理器挪到启动时（c56255c + 8924384）
  - 问题：CreateArticle/CreateExchangeRate/Register 每次请求都跑一次 AutoMigrate（=
    每请求执行 DDL）。
  - 改法：统一到 config/db.go 的 initDB()，启动时一次性迁移 Article/User/ExchangeRate
    三张表（顺带补上了 initDB 里原来漏掉的 ExchangeRate）。
  - 讲点：AutoMigrate 只该跑一次；启动路径 vs 请求路径的职责划分。

  四、安全

  7. JWT 密钥去硬编码（d8dab15）
  - 问题：密钥写死 []byte("secret")，代码泄露即可伪造 token，且无法轮换。
  - 改法：改读 config.AppConfig.JWT.Secret（来自 config.yml），并支持 JWT_SECRET
    环境变量覆盖，实现每环境不同密钥；密钥换成 openssl rand -hex 32 的随机串。
  - 讲点：hardcoded secret 的危害、配置/环境变量分层、密钥轮换的代价（换密钥 = 旧 token
    全部失效）。

  五、工程化

  8. 配置模板补占位符（2a842ec）
  - 给 config.yml.example 补上 jwt.secret 占位符。真实
    config.yml（含密钥和数据库密码）由仓库根 .gitignore 排除，不入库。
  - 讲点：真实凭据隔离、example 模板的协作习惯。

  六、运行时 Bug

  9. CORS 预检 403（598824f）
  - 问题：白名单只写 localhost:5173，用 127.0.0.1:5173 打开前端时 Origin
    不匹配，OPTIONS 预检被 403 拦下，表现为「注册/登录失败」。
  - 改法：白名单补 http://127.0.0.1:5173。
  - 讲点：CORS 预检机制、Origin 匹配的坑。
