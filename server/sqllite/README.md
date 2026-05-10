podman build -t test-api .
wildan-fq@nusantara:~$ mkdir -p ~/podman-data
wildan-fq@nusantara:~$ podman run -d \
  --name api-notes-test \
  -p 8080:8080 \
  -v ~/podman-data:/root/data:Z \
  test-api
b9425818ba43a917110ec056254045ee5210575a910145b1e9fdaafa64e8f73a
wildan-fq@nusantara:~$ podman ps
CONTAINER ID  IMAGE                      COMMAND       CREATED        STATUS        PORTS                   NAMES
b9425818ba43  localhost/test-api:latest  ./api-binary  4 seconds ago  Up 5 seconds  0.0.0.0:8080->8080/tcp  api-notes-test
wildan-fq@nusantara:~$
