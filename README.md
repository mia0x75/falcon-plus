[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
[![Badge](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu/#/zh_CN)

# Falcon+
```txt
    ___       ___       ___       ___       ___       ___    
   /\  \     /\  \     /\__\     /\  \     /\  \     /\__\   
  /  \  \   /  \  \   / /  /    /  \  \   /  \  \   / | _|_  
 /  \ \__\ /  \ \__\ / /__/    / /\ \__\ / /\ \__\ /  |/\__\ 
 \/\ \/__/ \/\  /  / \ \  \    \ \ \/__/ \ \/ /  / \/|  /  / 
    \/__/    / /  /   \ \__\    \ \__\    \  /  /    | /  /  
             \/__/     \/__/     \/__/     \/__/     \/__/   
```

# Enhancements

- 合并数据库到一个单独的库
- 统一配置文件存储位置
- 统一日志文件存储位置
- 命名一致性调整(例如: CreateUser, UserUpdate => verb + noun)
- 使用go.mod替换vendor方式
- 用户密码到md5+salt散列算法替换为bcrypt方式
- 所有的日志模块更改为logrus
- 移除viper读取配置到方式 统一使用读取本地JSON文件到方式
- 模块API引入缓存机制
- 代码一致性调整
- 统一组织所有模块到端口
- 统一替换time.Sleep为time.Tick
- 增加指标proc.num
- 移除指标agent.alive 该指标可用proc.num/name=falcon-agent替代
- 移除指标mysql.alive 该指标可用proc.num/name=mysqld替代
- 补齐API模块的自监控代码
- 性能计数器 (gateway.runtime/gateway.debug/graph.runtime/graph.debug)
- 统一使用日志级别设置log.level = "debug"替代"debug: true"设置
- 对agent模块增加webroot功能
- 更改min-step到5秒(rrd related)
- 增加自监控模块(exporter)
- 支持多心跳服务器(HBS)设置
- 移除GPU相关到监控指标
- 移除告警模块到回掉
- 推送原始JSON数据到告警模块
- 增强心跳服务器(HBS)
- 调制告警模块和告警判定模块中Redis的链接: [scheme](https://www.iana.org/assignments/uri-schemes/prov/redis)
- 代码重构
- 潜在的问题修复
- 合并Gateway的代码到Transfer
- 为Gateway, Exporter, Graph模块增加统计指标
- 增加自动更新模块

# TODO
- [TODO] exporter从hbs获取所有非维护状态的主机，并监控它们的健康状态，不用手工维护cfg.json
- [TODO] 增加文件系统文件变更监控，对重要文件的变更发出告警
- [TODO] 告警配置支持表达式，而不仅仅是常量，比如对于mysql.Thread_running > 2 * cpu.num + 2发出告警
- [TODO] API查询历史数据时，对于当前数据合并缓存内容，确保数据展现的及时性
- [TODO] mysql监控插件
- [TODO] redis监控插件
- [TODO] mongodb监控插件
- [TODO] 整合滴滴的日志内容监控

# 0.4.0改动

alarm.json
```
[+] redis.waittime_timeout 不可以为空
[>] redis.queue.* => queue.*
```


judge.json
```
[>] alarm.redis.dsn => alarm.redis.addr
[>] alarm.redis.connect_timeout => unit second
[>] alarm.redis.read_timeout => unit second
[>] alarm.redis.write_timeout => unit second
[+] alarm.redis.waittime_timeout 不可以为空
[-] alarm.enable 必须enable
[>] alarm.redis.* => redis.*
```

transfer.json
```
[+] ignore.*
```

# Documentations

- [Usage](http://book.open-falcon.com)
- [Open-Falcon API](http://open-falcon.com/falcon-plus)

# Prerequisite

- git >= 1.7.5
- go >= 1.6
- upx

# Getting Started

## Build from source
**before start, please make sure you prepared this:**

```
yum install -y redis
yum install -y mysql-server

```

*NOTE: be sure to check redis and mysql-server have successfully started.*

And then

```
# Please make sure that you have set `$GOPATH` and `$GOROOT` correctly.
# If you have not golang in your host, please follow [https://golang.org/doc/install] to install golang.

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus.git

```

**And do not forget to init the database first (if you have not loaded the database schema before)**

```
cd $GOPATH/src/github.com/open-falcon/falcon-plus/scripts/mysql/db_schema/
mysql -h 127.0.0.1 -u root -p < 1_dashboard-db-schema.sql
```

# Compilation

```
cd $GOPATH/src/github.com/open-falcon/falcon-plus/

# make all modules
make all

# make specified module
make agent

# pack all modules
make pack
```

* *after `make pack` you will got `open-falcon-vx.x.x.tar.gz`*
* *if you want to edit configure file for each module, you can edit `config/xxx.json` before you do `make pack`*

#  Unpack and Decompose

```
export WorkDir="$HOME/open-falcon"
mkdir -p $WorkDir
tar -xzvf open-falcon-vx.x.x.tar.gz -C $WorkDir
cd $WorkDir
```

# Start all modules in single host
```
cd $WorkDir
./open-falcon start

# check modules status
./open-falcon check

```

# Run More Open-Falcon Commands

for example:

```
# ./open-falcon [start|stop|restart|check|monitor|reload] module
./open-falcon start agent

./open-falcon check
        falcon-graph         UP           53007
          falcon-hbs         UP           53014
        falcon-judge         UP           53020
     falcon-transfer         UP           53026
       falcon-nodata         UP           53032
   falcon-aggregator         UP           53038
        falcon-agent         UP           53044
      falcon-gateway         UP           53050
          falcon-api         UP           53056
        falcon-alarm         UP           53063
```

* For debugging , You can check `$WorkDir/$moduleName/logs/xxx.log`

# Install Frontend Dashboard
- Follow [this](https://github.com/open-falcon/dashboard).

**NOTE: if you want to use grafana as the dashboard, please check [this](https://github.com/open-falcon/grafana-openfalcon-datasource).**

# Package Release

```
make clean all pack
```

# API Standard
- [API Standard](https://github.com/open-falcon/falcon-plus/blob/master/api-standard.md)


# Q&A

- Any issue or question is welcome, Please feel free to open [github issues](https://github.com/open-falcon/falcon-plus/issues) :)
- [FAQ](http://book.open-falcon.com/zh_0_2/faq/)

