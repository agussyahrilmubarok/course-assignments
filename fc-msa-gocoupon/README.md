# Go Coupon MSA


## Run Application

1. Run docker-compose

```bash
make compose/local/up
# OR
make compose/up
```

2. Run app

```bash
cd app && \
make migrater/up && \
make api/dev
```

3. Run Testing

```bash
cd k6/$(version) && \
k6 run $(test_file)
```

4. Explore 

```bash
# Postgres
http://localhost:5050

# Redis
http://localhost:5051

# Kafka
http://localhost:5052

# Grafana
http://localhost:3000

# Prometheus
http://localhost:9090

# Zipkin
http://localhost:9411
```