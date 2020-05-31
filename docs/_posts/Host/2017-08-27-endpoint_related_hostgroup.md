---
category: Host
apiurl: '/api/v1/endpoint/#{name}/hostgroup'
title: "Get related HostGorup of Endpoint"
type: 'GET'
sample_doc: 'host.html'
layout: default
---
 * [Session](#/authentication) Required
* ex. /api/v1/endpoint/test1/hostgroup
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
