# 快速上手
每个租户都有自己的用户目录，管理所有第三方登陆。具体可以参考 *oauth2* 协议

1. 创建一个`client`。并创建`redirect_uri`
2. 登陆用户时，请求`/oauth2/auth`接口获`authorization_code`得，如果没有登陆，会跳转到登陆界面。 如果自定义登陆页面，可以用登陆页面直接请求`/login`接口进行用户名密码登陆，也可以请求`/login/{provider}`使用对应的第三方登陆方式。具体支持哪些登陆方式可以使用`/login/providers`进行查询，然后使用`/login/providers/{provider}`查询具体登陆配置。
3. 使用上一步获得的`authorization_code`，请求`/oauth2/token`获得`access_token`