1.安装


2. 拷贝 webhooks 到 `/opt` 目录下

3. 创建配置文件

   ```
   touch /etc/supervisord.d/webhook.ini
   ```

   添加以下内容

   ```
   [program:webhooks]
   command=/opt/webhooks/webhooks -c /opt/webhooks/log/webhooks.yaml  ;应用入口
   user=root
   directory=/opt/webhooks                         ;web目录
   startsecs=5                                      ;启动时间
   stopwaitsecs=0                                  ;终止等待时间
   startretries = 3
   autostart=true                                   ;是否自动启动
   autorestart=true                                 ;是否自动重启
   redirect_stderr=true                             ;错误日志输出到标准日志
   stdout_logfile=/dev/null                      ;标准日志不输出
   stdout_logfile_maxbytes=10MB                        ;标准日志大小
   ```

   

4. 开启数据库支持

  打开yaml 文件
  ```
  mysql:
  enable: true                   # 改为true 并配置数据库账号密码

  
  ```
   

5. 使 supervisor 自启动

   ```
   systemctl enable supervisord
   ```

   测试 

   ```
   curl -s http://localhost:21332/hookssync/38daa500-55a8-11e8-acd6-704d7b885ead
   ```

注意事项

1.CentOS 防火墙记得放行21332端口或者直接关闭防火墙`systemctl stop firewalld;systemctl disable firewalld`

dev 
```
set GO111MODULE=on

go mod init  projectName
go mod tidy
```

go mod help查看帮助
go mod init<项目模块名称>初始化模块，会在项目根目录下生成 go.mod文件。

go mod tidy根据go.mod文件来处理依赖关系。

go mod vendor将依赖包复制到项目下的 vendor目录。建议一些使用了被墙包的话可以这么处理，方便用户快速使用命令go build -mod=vendor编译

go list -m all显示依赖关系。go list -m -json all显示详细依赖关系。

go mod download <path@version>下载依赖。参数<path@version>是非必写的，path是包的路径，version是包的版本。
