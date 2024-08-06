docker build -t cs-server-manager --progress=plain .
docker run -it --rm --name cs-server-manager --mount type=bind,source=//home/desk/programming/files/cs,destination=/data -p 8080:80 -p 27015:27015 cs-server-manager 
