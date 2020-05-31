---
category: HostGroup
apiurl: '/api/v1/hostgroup/template'
title: "Bind A Template to HostGroup"
type: 'POST'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```{
  "template_id": 5,
  "group_id": 3
}```

### Response

```Status: 200```
```{"group_id":3,"template_id":5,"bind_user":2}```
