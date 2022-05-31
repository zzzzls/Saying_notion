# 简介

> 使用 Notion 官方 API

平时总摘抄句子到 Notion 中，却很少回顾

该程序随机从 Notion 数据库中返回一条句子，并展示在首页。



# 使用示例

1. 前往 Notion API 站点生成 Token

2. 前往对应数据库及展示页面授权 Token

3. 执行命令：

   ```shell
   # windows
   saying.exe -toekn <notion token> -did <Database Id> -bid <Block Id>
   
   # linux
   ./saying -toekn <notion token> -did <Database Id> -bid <Block Id>
   ```

​		

​		
