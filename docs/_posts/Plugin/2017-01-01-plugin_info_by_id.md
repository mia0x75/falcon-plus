---
category: Plugin
apiurl: '/api/v1/hostgroup/#{id}/plugins'
title: "Get Plugin List of HostGroup"
type: 'GET'
sample_doc: 'plugin.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup/343/plugins
* group_id: hostgroup id

### Response

```Status: 200```
```[
  {
    "id": 1499,
    "group_id": 343,
    "dir": "testpath",
    "create_user": "root"
  },
  {
    "id": 1501,
    "group_id": 343,
    "dir": "testpath",
    "create_user": "root"
  }
]```
