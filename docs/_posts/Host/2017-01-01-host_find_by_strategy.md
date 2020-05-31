---
category: Host
apiurl: '/api/v1/host/find_by_strategy'
title: "Find host by strategy"
type: 'POST'
sample_doc: 'host.html'
layout: default
---
 * [Session](#/authentication) Required
 ### Request
```
{
	"metric": "agent.alive"
}
```
 ### Response
 ```Status: 200```
```
[
    {
        "strategy": {
            "id": 13,
            "metric": "agent.alive",
            "tags": "",
            "max_step": 1,
            "priority": 2,
            "func": "all(#10)",
            "op": "<",
            "right_value": "0",
            "note": "machine is down",
            "run_begin": "",
            "run_end": "",
            "template_id": 14
        },
        "hosts": [
            "foo1.bar.com",
            "foo2.bar.com",
            "foo3.bar.com",
        ]
    },
    ...
]
```
