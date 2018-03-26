# Generate Cassandra Trigger Database Instance

To build docker image. By default only the 9042 port is exposed, 
so you will have to modify the docker file if you want to do any clustering

```
docker build -t <MyDockerImageName> .
```

To use the image available from my docker cloud account

```
docker pull asinha94/seng468_cassandra_transaction
```

To use the image

```
docker run asinha94/seng468_cassandra_transaction
```

This container generates a lot of output, so you may want to send it to the background and ignore it

```
docker run -d asinha94/seng468_cassandra_transaction
```

