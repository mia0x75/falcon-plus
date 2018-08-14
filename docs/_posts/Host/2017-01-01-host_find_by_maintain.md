---
category: Host
apiurl: '/api/v1/hosts/maintain'
title: "Find host in maintain"
type: 'GET'
sample_doc: 'host.html'
layout: default
---
 * [Session](#/authentication) Required
 ### Response
 ```Status: 200```
```
[
    {
        "id": 7101,
        "hostname": "foo1.bar.com",
        "ip": "10.1.1.1",
        "agent_version": "6.0.1",
        "plugin_version": "plugin not enabled",
        "maintain_begin": 1502781600,
        "maintain_end": 1508052000
    },
    {
        "id": 47313,
        "hostname": "foo2.bar.com",
        "ip": "10.1.1.2",
        "agent_version": "",
        "plugin_version": "",
        "maintain_begin": 1502851980,
        "maintain_end": 1502862780
    }
]
```
