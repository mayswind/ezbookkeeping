# ezBookkeeping 部署说明（deploy/）

> 这套文件用来在服务器上部署 ezBookkeeping：MySQL + 应用 + Caddy（自动 HTTPS）。
> 配置的「真身」在本仓库 `deploy/`（有版本管理）；服务器上只放一份副本 + 一个 `.env`。
> 应用镜像由 GitHub Actions（CI）构建并推到 Docker Hub，服务器只拉成品镜像，**不在服务器编译源码**。

## 文件说明

| 文件 | 作用 | 进 git | 服务器上 |
|---|---|---|---|
| `docker-compose.yml` | 编排 MySQL + 应用 + Caddy，敏感值用 `${变量}` 引用 `.env` | ✅ | 整文件覆盖 |
| `Caddyfile` | 反向代理 + 自动 HTTPS；域名从 `{$EBK_SERVER_DOMAIN}` 注入 | ✅ | 整文件覆盖 |
| `init.sql` | 仅 MySQL **首次初始化**时建 stg/prod 两库并授权 | ✅ | 整文件覆盖（已初始化则是空操作） |
| `.env.example` | 变量清单模板，不含真实值 | ✅ | 作填写参考 |
| `.env` | 真实密码/密钥/域名，`.gitignore` 已忽略 | ❌ | **只在服务器手建手改** |

## 两条铁律

1. **四个文件必须在同一目录**：compose 用相对路径挂载 `./Caddyfile`、`./init.sql`，`.env` 也从「执行命令的目录」自动读取。缺一个或不在同目录，`up -d` 会报错。
2. **`.env` 里的数据库账号密码必须和 MySQL 数据卷初始化时的值一致**：密码在首次初始化时就固化进 `mysql_data` 卷了，填错应用连不上库（日志报认证失败）。想换密码是独立操作（见末尾）。

> 目录路径本身随意（不限 `~/apps/ezbookkeeping`）：compose 里钉了 `name: ezbookkeeping`，数据卷绑定的是项目名而非目录名，换目录跑数据照样在。

---

## 一、首次部署（从零到上线，一次性）

> 换台新服务器重来，照这个走即可。已部署过则跳到「二、日常更新」。

### 1. 服务器准备
```bash
# 安装 Docker + Compose（Ubuntu 官方一键脚本）
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER        # 免 sudo；执行后必须重新登录 SSH 才生效
docker version && docker compose version
```

### 2. 放置部署文件
在服务器建个目录（如 `~/apps/ezbookkeeping`），把本仓库 `deploy/` 下的
`docker-compose.yml`、`Caddyfile`、`init.sql` 三个文件整文件贴进去。
YAML 缩进敏感，用带单引号的 heredoc 最稳（`'EOF'` 能让 `${...}` 原样保留、不被 shell 提前展开）：
```bash
mkdir -p ~/apps/ezbookkeeping && cd ~/apps/ezbookkeeping
cat > docker-compose.yml <<'EOF'
# ……粘贴本仓库 deploy/docker-compose.yml 的全部内容……
EOF
cat > Caddyfile <<'EOF'
# ……粘贴 Caddyfile 全部内容……
EOF
cat > init.sql <<'EOF'
# ……粘贴 init.sql 全部内容……
EOF
```

### 3. 建 `.env` 填真实值
```bash
cp .env.example .env     # 若没贴 .env.example，就直接 nano .env 照模板写
nano .env                # 填密码/密钥/域名；EBK_SECURITY_SECRET_KEY 用 openssl rand -base64 32 生成
chmod 600 .env           # 只有自己可读
```

### 4. 云平台两件事（容易忘）
- **防火墙**放行入站 **22 / 80 / 443**（Caddy 申请证书 80、443 都要，少一个证书下不来）。
- **DNS** 给域名加 A 记录指向服务器公网 IP，等生效（`dig +short 你的域名`）。

### 5. 校验并启动
```bash
docker compose config | grep -i 'image\|passwd\|domain'   # 确认没有残留 ${...}
docker compose pull            # 拉镜像
docker compose up -d
docker compose ps              # mysql 应 healthy，app/caddy 应 Up
docker compose logs -f caddy   # 期望看到 certificate obtained successfully
```

### 6. 验证
浏览器打开 `https://你的域名`，看到带锁的页面即成功。**第一个注册的账号是管理员。**
服务器本地自测（正确 SNI）：
```bash
curl -Ik https://你的域名 --resolve 你的域名:443:127.0.0.1
```

---

## 二、日常更新

### A. 更新应用（改了源码、CI 已重建镜像）
```bash
cd ~/apps/ezbookkeeping
docker compose pull            # 拉 CI 刚构建的新镜像
docker compose up -d           # 用新镜像重启（数据卷不变，数据保留）
docker compose logs -f ezbookkeeping
# 确认跑的是自己的镜像：
docker inspect --format '{{.Config.Image}}' ezbookkeeping
```

### B. 改了 compose / Caddyfile（少见）
笔记本改 `deploy/` 里的文件 → push 留版本 → 服务器把对应文件**整文件重贴**（同第一步 heredoc）→
```bash
docker compose up -d           # 变更生效
```

### C. 改 `.env`（最常见，不碰 git、不等 CI）
```bash
cd ~/apps/ezbookkeeping
nano .env
docker compose up -d           # 立即生效
```

### 常用排错
```bash
docker compose ps                         # 谁挂了
docker compose logs --since=15m ezbookkeeping | grep -i error
docker compose logs -f caddy              # 域名/HTTPS/反代
docker compose logs mysql                 # 连不上库
docker compose exec ezbookkeeping sh      # 进容器内手查
```

---

## 三、切换 stg / prod

只改 `.env` 里一行 `EBK_DATABASE_NAME`（`ezbookkeeping_stg` ↔ `ezbookkeeping_prod`），然后 `docker compose up -d`。
两个库在首次初始化时已由 `init.sql` 一起建好。

## 四、切换到自己的镜像（方式 3 → 方式 4）

`.env` 里 `EBK_IMAGE` 默认已是 `gaofei211/bookkeeping:latest`（自己 CI 的镜像）。
想临时验证作者镜像改成 `mayswind/ezbookkeeping:latest` 即可。改完 `docker compose pull && up -d`。

## 五、邮件（SMTP）与忘记密码

`EBK_MAIL_*` 配好后，登录页「忘记密码」才真正能发重置邮件（`enable_forget_password` 默认即 true）。
用个人 QQ 邮箱发信：`mail.qq.com` 设置→账号 开启 IMAP/SMTP，发短信后拿到 16 位**授权码**（不是登录密码）。
把 `.env` 里 `EBK_MAIL_SMTP_*` / `EBK_MAIL_FROM_ADDRESS` 填好（见 `.env.example`），`docker compose up -d`，然后测试：
```bash
docker compose exec ezbookkeeping /ezbookkeeping/ezbookkeeping utility send-test-mail --to 你的邮箱
```
- 端口用 `465`（云服务器一般封出站 25）。QQ 要求发件人地址 = 认证邮箱，`FROM_ADDRESS` 必须填同一个 QQ 邮箱。
- 没配邮件又想临时重置某用户密码：
  `docker compose exec ezbookkeeping /ezbookkeeping/ezbookkeeping user-data user-modify-password --username <用户名> --password <新密码>`

## 六、备份与改密码（提醒，别在日常更新里顺手做）

```bash
# 备份（导出当前库）
docker exec ezbookkeeping-mysql mysqldump -uroot -p<root密码> ezbookkeeping_stg > backup_stg_$(date +%F).sql
```
- **改数据库密码**：要同时改 MySQL 里用户的密码 **和** `.env`，两边一致才连得上，不能只改一边。
- **改 `EBK_SECURITY_SECRET_KEY`**：改了所有已登录用户的 token 失效、需重新登录（数据不受影响）。

---

### 备忘
- `init.sql` 只在 MySQL **数据卷为空**时执行一次；库已建好后再贴也不会重跑。补建库要用
  `docker exec ... mysql -e "CREATE DATABASE ..."`。
- Caddy 自动 HTTPS 需要 80 + 443 都放行；证书 90 天自动续期。
- MySQL 不对宿主机映射端口（compose 里 mysql 无 `ports`），只在 Docker 内网可达，公网连不到。
