## 🐳 Hello Docker — Minimal Go Application

A lightweight Go application running inside a Docker container.
This image prints a simple greeting message and shows its build version.

### Features

* Based on Go static binary

### Usage

```bash
docker pull agussyahrilmubarok/hello-docker:latest
docker run --rm agussyahrilmubarok/hello-docker:latest
```

Example output:

```
Hello from Docker! This Go application is running inside a container.
Version: 1.0.0
```

Use a specific version for stable deployment:

```bash
docker pull agussyahrilmubarok/hello-docker:latest
```

Run as container:

```bash
docker run --rm agussyahrilmubarok/hello-docker:latest
```

### Recommended Use Cases

* Testing Docker environments
