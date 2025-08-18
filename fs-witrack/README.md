# WiTrack

WiTrack is a web-based ticket tracking system designed to streamline the management of Wi-Fi installation services. It allows businesses to efficiently track, assign, and resolve service tickets, ensuring timely support and improved customer satisfaction. The application features a Golang backend for robust, high-performance processing, and a Vue.js frontend for a responsive and intuitive user interface.

## Features

* **Ticket Management**: Create, update, assign, and close Wi-Fi installation tickets.

## Technology Stack

* [Java](https://www.java.com/)
* [Spring](https://spring.io/)
* [MongoDB](https://www.mongodb.com/)
* [Testcontainers]
* [Vue.js](https://vuejs.org/)
* [Tailwind CSS](https://tailwindcss.com/)
* [JWT](https://jwt.io/)

## Explore

```bash
# Run Infrastructure
docker compose up -d --build

# Run API
go mod tidy
go run cmd/main.go

# Run Frontend
pnpm install
pnpm run dev

# Run Production Level
docker compose -f docker-compose-prod.yml up -d --build
```

## References

* [Java Documentation](https://www.java.com/en/docs/)
* [Spring Framework](https://spring.io/)
* [MongoDB Documentation](https://www.mongodb.com/docs/)
* [Vue.js Guide](https://vuejs.org/guide/introduction.html)
* [Tailwind CSS Documentation](https://tailwindcss.com/docs)
* [JWT Introduction](https://jwt.io/introduction/)