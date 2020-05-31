---
category: Plugin
apiurl: '/api/v1/plugin'
title: "Create A Plugin of HostGroup"
type: 'POST'
sample_doc: 'plugin.html'
layout: default
---

* [Session](#/authentication) Required
* group_id: hostgroup id

### Request

```{
  "hostgroup_id": 343,
  "dir_path": "testpath"
}```

### Response

```Status: 200```
```{
  "id": 1501,
  "group_id": 343,
  "dir": "testpath",
  "create_user": "root"
}```
