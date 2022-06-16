# 简介

changeDetector, 通过cron定时任务访问指定页面，检测到页面变更时，会通过已配置的联系方式通知用户。可用用来追踪游戏、应用更新，用户文章更新等。

# 截图



# 使用方式

1. 添加页面解析规则

   页面解析规则说明了从网页中捕获信息的方式，使用CSS选择器，该规则集定义了以下5个字段

   - name：规则的名称
   - md5：通过该选择器获取到一段文本，当该文本变更时，表面页面发生了更新
   - title：通过该选择器获取到一段文本，用于生成通知的标题
   - body：通过该选择器获取到一段文本，用于生成通知的主要内容
   - link：通过该选择器获取到一个a标签，a标签的href属性为通知提供一个链接。为空时，取解析页面链接

2. 添加要追踪的页面

   - name：为该条目起一个名字
   - type：规则的名称，说明页面解析方式
   - url：追踪页面的链接

3. 添加通知方式

   目前只支持通过邮件通知，需要填写4个字段

   - smtpUser：邮箱账号，例如：example@outlook.com
   - smtpPswd：邮箱密码
   - smtpServer：邮箱服务器地址，例如：smtp-mail.outlook.com
   - smtpPort：邮箱服务器端口，例如：587

4. 设置检查更新的周期

   用的[cron](https://github.com/robfig/cron)库，支持其内置的语法，例如

   - 标准crontab语法： * */12 * * * 表示每12h执行一次，30 * * * * 表示每个小时的第30分时执行一次，* * 2,4,6,8 * * 表示每月2、4、6、8日执行
   - @yearly、@monthly、@weekly、@daily、@hourly 语法，顾名思义
   - @every 1h、@every 5m、@every 1h30m 语法，顾名思义



# 编译

1. 安装Go开发环境

2. 下载此仓库

   ```bash
   git clone https://github.com/GeminiT369/changeDetector
   ```

3. 进入项目

   ```bash
   cd changeDetector
   ```

4. 编译

   ```bash
   go build .
   ```

   

5. 