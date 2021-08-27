```json
{
  "groups": [
    {
      "group": "rest-api",
      "containers": ["instance-1", "instance-2", "instance-3"],
      "image_tag": "2.6.0",
      "endpoints": [
        {
          "location": "https://node1.test:8888",
          "node_name": "node1"
        },
        {
          "location": "https://node2.test:8888",
          "node_name": "node2"
        }
      ]
    }
  ]
}
```

```json
{
  "node_name": "personal",
  "is_online": true,
  "containers": [
    {
      "container_name": "portainer",
      "container_id": "311c42383ec3",
      "image": "portainer/portainer-ce:2.6.2",
      "status": "Up 8 hours"
    }
  ]
}
```