-- 仅在 MySQL 数据卷首次初始化（空卷）时执行一次。
-- 数据卷非空后不再运行；以后补建库要用 docker exec ... mysql -e "CREATE DATABASE ..."。
-- 这里只建库 + 授权，不含任何密码（用户与密码由 compose 的环境变量创建）。
-- 注意：下面授权语句里的用户名 'ezbookkeeping' 是写死的，静态 SQL 读不到 .env 变量，
--       所以它必须和 .env 里的 EBK_DATABASE_USER 保持一致；若改了用户名，这里也要同步改。
CREATE DATABASE IF NOT EXISTS ezbookkeeping_stg  CHARACTER SET utf8mb4;
CREATE DATABASE IF NOT EXISTS ezbookkeeping_prod CHARACTER SET utf8mb4;
GRANT ALL PRIVILEGES ON ezbookkeeping_stg.*  TO 'ezbookkeeping'@'%';
GRANT ALL PRIVILEGES ON ezbookkeeping_prod.* TO 'ezbookkeeping'@'%';
FLUSH PRIVILEGES;
