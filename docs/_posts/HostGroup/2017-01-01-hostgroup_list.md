---
category: HostGroup
apiurl: '/api/v1/hostgroup'
title: "HostGroup List"
type: 'GET'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required

### Response

```Status: 200```
```[
  {
    "id": 3,
    "name": "docker-A",
    "create_user": "user-A"
  },
  {
    "id": 5,
    "name": "docker-T",
    "create_user": "user-B"
  },
  {
    "id": 8,
    "name": "docker-F",
    "create_user": "root"
  }
]```
