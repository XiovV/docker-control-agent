# Setup

## Docker install
### Step 1: Run the agent
```shell
docker run -d -p 8080:8080 --name=dokkup-agent \
    --restart=always \
    -v /var/run/docker.sock:/var/run/docker.sock \
    xiovv/dokkup-agent:latest
```

### Step 2: Get your API key
```shell
docker logs dokkup-agent
```

Output:
```shell
Your new api key is: DSK7D4TL5LIJT5R5LVCUCOBHQ4
Successfully loaded config
agent is listening on :8080
```