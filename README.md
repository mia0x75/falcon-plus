[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
[![Badge](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu/#/zh_CN)

# Falcon+

![Open-Falcon](./logo.png)

# 系统改进

- Add updater module
- Rename falcon_portal to portal
- Potential bug fix
- Put all config files to /etc/fpm
- Change all output log files to /var/log/fpm
- Naming consistency
- Builtin mysql monitor
- Builtin redis monitor
- Update vendor
- Replace md5+salt with bcrypt
- Change log with logrus in all modules
- Replace viper with local cfg code
- Code format and style consistency
- Re-organize port usage (xxx0 - RPC / xxx1 - HTTP / xxx2 - Sockets)
- Replace time.Sleep with time.Tick
- Add metric proc.num
- Remove agent.alive, agent.alive replace with proc.num/name=falcon-agent
- Remove mysql.alive, see above
- Add health http method for api module
- PFC (Performance Counter) (gateway.runtime/gateway.debug/graph.runtime/graph.debug)
- Support interval settings for echo func
- Use log.level = "debug" for debug mode, remove "debug: true"
- Add webroot for agent
- Change min-step to 5 sec (rrd related)
- Add exporter module
- Add health monitor into exporter module
- Support multiple hbs in agent module
- Remove GPU related metrics
- Save statistic metrics from /statistics/all (gateway.stats/exporter.stats/graph.stats)
- Code refactor


transfer.json
[+] ignore.*

# Documentations

- [Usage](http://book.open-falcon.com)
- [Open-Falcon API](http://open-falcon.com/falcon-plus)

# Prerequisite

- Git >= 1.7.5
- Go >= 1.6

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
mysql -h 127.0.0.1 -u root -p < 1_uic-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 2_portal-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 3_dashboard-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 4_graph-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 5_alarms-db-schema.sql
```

**NOTE: if you are upgrading from v0.1 to current version v0.2.0,then**. [More upgrading instruction](http://www.jianshu.com/p/6fb2c2b4d030)

    mysql -h 127.0.0.1 -u root -p < 5_alarms-db-schema.sql

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

# Package Management

We use govendor to manage the golang packages. Please install `govendor` before compilation.

    go get -u github.com/kardianos/govendor

Most depended packages are saved under `./vendor` dir. If you want to add or update a package, just run `govendor fetch xxxx@commitID` or `govendor fetch xxxx@v1.x.x`, then you will find the package have been placed in `./vendor` correctly.

# Package Release

```
make clean all pack
```

# API Standard
- [API Standard](https://github.com/open-falcon/falcon-plus/blob/master/api-standard.md)


# Q&A

- Any issue or question is welcome, Please feel free to open [github issues](https://github.com/open-falcon/falcon-plus/issues) :)
- [FAQ](http://book.open-falcon.com/zh_0_2/faq/)

