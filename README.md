# 用户与权限管理系统

本系统主要有两个功能
1. 管理用户身份信息
2. ABAC方式管理资源权限信息
3. 资源组权限管理

# 租户
本系统为多租户部署，请求url为`https://{host}/[{tenant}/]/*`，`{tenant}`为可选项，会按照如下顺序匹配租户：
1. 从用tenant参数匹配租户
2. 从host匹配租户
3. 默认default租户


# cmd

- 生成文档

```bash
swag init -o backend/docs
```

- docker部署

```bash
docker-compose up -d
```
