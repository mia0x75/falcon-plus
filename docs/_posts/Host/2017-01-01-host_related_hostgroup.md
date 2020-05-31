---
category: Host
apiurl: '/api/v1/host/#{id}/hostgroup'
title: "Get related HostGorup of Host"
type: 'GET'
sample_doc: 'host.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/host/1647/hostgroup
* name: hostgroup name

### Response

```Status: 200```
```[
  {
    "id": 78,
    "name": "tplB",
    "create_user": "userA"
  },
  {
    "id": 145,
    "name": "Owl_Default_Group",
    "create_user": "userA"
  }
]```
