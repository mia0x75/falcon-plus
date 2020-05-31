#!/bin/bash

confs=(
	'%%REDIS%%=redis:\/\/127.0.0.1:6379\/0'
	'%%MYSQL%%=root:root@tcp(127.0.0.1:3306)'

	'%%API_DEFAULT_TOKEN%%=TOKEN'
	'%%API_RPC%%=127.0.0.1:6000'
	'%%API_HTTP%%=127.0.0.1:6008'

	'%%EXPORTER_HTTP%%=127.0.0.1:6018'

	'%%UPDATER_HTTP%%=127.0.0.1:6028'

	'%%HBS_RPC%%=127.0.0.1:6030'
	'%%HBS_HTTP%%=127.0.0.1:6038'

	'%%AGGREGATOR_HTTP%%=127.0.0.1:6048'

	'%%ALARM_HTTP%%=127.0.0.1:6058'

	'%%TRANSFER_RPC%%=127.0.0.1:6060'
	'%%TRANSFER_HTTP%%=127.0.0.1:6068'
	'%%TRANSFER_SOCKET%%=127.0.0.1:6062'

	'%%GRAPH_RPC%%=127.0.0.1:6070'
	'%%GRAPH_HTTP%%=127.0.0.1:6078'

	'%%JUDGE_RPC%%=127.0.0.1:6080'
	'%%JUDGE_HTTP%%=127.0.0.1:6088'

	'%%NODATA_HTTP%%=127.0.0.1:6098'

	'%%GATEWAY_RPC%%=127.0.0.1:6100'
	'%%GATEWAY_HTTP%%=127.0.0.1:6108'
	'%%GATEWAY_SOCKET%%=127.0.0.1:6102'

	'%%AGENT_HTTP%%=127.0.0.1:6818'
)

configurer() {
	for i in "${confs[@]}"
	do
		search="${i%%=*}"
		replace="${i##*=}"

		find ./out/*/*.json -type f -exec sed -i -e "s/${search}/${replace}/g" {} \;
	done
}
configurer
