## Starting project

1. `docker-compose up -d`
2. `make serve`

## Debezium (Event driven)

1. Start Kafka Connector: (TODO: Add to docker compose)
   `docker run --name connect --network easyisa_default -p 8083:8083 -e BOOTSTRAP_SERVERS=easyisa_kafka_1:9092 -e GROUP_ID=1 -e CONFIG_STORAGE_TOPIC=easyisa_connect_configs -e OFFSET_STORAGE_TOPIC=easyisa_connect_offsets -e STATUS_STORAGE_TOPIC=easyisa_connect_statuses --link easyisa_zookeeper_1:zookeeper --link easyisa_kafka_1:kafka --link easyisa_postgres_1:postgres debezium/connect:0.10`

2. Register the default topic using cadence-cli
   `docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain easyisa-domain domain register`

3. Connect Prisma DB to Kafka (TODO: Add to Makefile)
   `curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" http://localhost:8083/connectors/ -d @connector.json`

4. Start sample consumer (For debugging)
   `docker run -it --name watcher --rm --network easyisa_default -e KAFKA_BROKER=easyisa_kafka_1:9092 -e ZOOKEEPER_CONNECT=easyisa_zookeeper_1:2181 --link easyisa_zookeeper_1:zookeeper --link easyisa_kafka_1:kafka debezium/kafka:0.10 watch-topic -a -k prisma.default_default.Post`
