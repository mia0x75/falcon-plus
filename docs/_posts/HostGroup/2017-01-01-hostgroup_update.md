---
category: HostGroup
apiurl: '/api/v1/hostgroup/#{id}/update'
title: "Update HostGroup info by id"
type: 'PUT'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup

### Request
```{
  "id" : 343,
  "name": "test1"
}```

### Response

```Status: 200```
```{
  "message":"hostgroup profile updated"
}```
