# 重点研发计划-课题二【2019YFB2101702】

课题任务书和相关文档见：`/doc`

## 代码

本仓库由以下两部分源码组成

1. 后端：`/backend`
2. 前端：`/fe`



## 静态资源

打包好的前端静态资源见本仓库的 **Releases**，可以从中下载**分布式公钥基础设施**的前端静态资源（是一个压缩包），静态资源有以下两种使用方式

**1. 与本地的后端调试【TODO】**

**2. 部署至服务器**

需要配合 Nginx 使用，配置文件在`/etc/nginx/sites-available/default`，参考配置如下

```nginx
server{
	location /cert {
		alias /home/zdyf/cert/;
		#指定主页
		index index.html;
		#自动跳转
		autoindex on;
	}
}
```

解释：

用户在浏览器访问 `<your-domain>/cert` 时，请求打到服务器上交由 Nginx 处理，Nginx 根据配置的内容返回服务器上 `/home/zdyf/cert/` 路径下的静态文件（`.html`、`.css`、`.js` 等）给浏览器，最后展示完整的前端页面。

所以**部署前端**实际上就是把本地打包好的（数据共享交换平台）静态资源上传至`/home/zdyf/cert/`目录下

最终服务器上的目录结构如下：

```
home/zdyf
|
+---cert
|   |   index.html
|   +---css
|   +---js
```



