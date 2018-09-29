#!/bin/bash

confs=(
    '%%REDIS%%=127.0.0.1:6379'
    '%%MYSQL%%=root:root@tcp(127.0.0.1:3306)'
    '%%API_DEFAULT_TOKEN%%=Y3!dHXtN]c>k@u4O.F,?L}/o=W*pgv-&'
    '%%API_HTTP%%=127.0.0.1:6001'
    '%%HBS_RPC%%=127.0.0.1:6030'
    '%%HBS_HTTP%%=127.0.0.1:6031'
    '%%AGGREGATOR_HTTP%%=127.0.0.1:6041'
    '%%TRANSFER_RPC%%=127.0.0.1:6060'
    '%%TRANSFER_HTTP%%=127.0.0.1:6061'
    '%%TRANSFER_SOCKET%%=127.0.0.1:6062'
    '%%GRAPH_RPC%%=127.0.0.1:6070'
    '%%GRAPH_HTTP%%=127.0.0.1:6071'
    '%%JUDGE_RPC%%=127.0.0.1:6080'
    '%%JUDGE_HTTP%%=127.0.0.1:6081'
    '%%NODATA_HTTP%%=127.0.0.1:6091'
    '%%ALARM_HTTP%%=127.0.0.1:7001'
    '%%GATEWAY_RPC%%=127.0.0.1:7010'
    '%%GATEWAY_HTTP%%=127.0.0.1:7011'
    '%%GATEWAY_SOCKET%%=127.0.0.1:7012'
    '%%AGENT_HTTP%%=127.0.0.1:7021'
 )

configurer() {
    for i in "${confs[@]}"
    do
        search="${i%%=*}"
        replace="${i##*=}"

        uname=`uname`
        if [ "$uname" == "Darwin" ] ; then
            # Note the "" and -e  after -i, needed in OS X
            find ./out/*/*.json -type f -exec sed -i .tpl -e "s/${search}/${replace}/g" {} \;
        else
            find ./out/*/*.json -type f -exec sed -i "s/${search}/${replace}/g" {} \;
        fi
    done
}
configurer
