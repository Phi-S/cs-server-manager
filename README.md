# cs-server-manager

```
docker run -it --rm --name cs-server-manager --mount type=bind,source=/cs-server-manager,destination=/data -e HTTP_PORT=8080 -e CS_PORT=27015 -e DATA_DIR=/data -p 8080:8080 -p 27015:27015 cs-server-manager
```