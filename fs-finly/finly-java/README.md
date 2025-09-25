# **Finly - Invoicing Web Application**

**Finly** is a modern invoicing application designed to help businesses easily generate and manage invoices online. It includes features like invoice creation, management, and data visualization through interactive charts.

This project has been built using **Java** with the **Spring** framework, **MongoDB**, and **Chart.js** for data visualization. The frontend is styled with **Tailwind CSS** and **DaisyUI**.

---

## **Features**

- Create and manage invoices online.
- View detailed invoice statistics.
- Interactive data visualization with charts.
- User-friendly and responsive design powered by Tailwind CSS and DaisyUI.

---

## **Requirements**

- **Java** (v21 or higher)
- **MongoDB** (v7 or higher)

---

## **Getting Started**

1. Install dependencies:

   ```bash
   mvn clean install verify
   ```

2. Copy the environment variables file:

   ```bash
   cp .env.sample .env
   ```

3. Ensure mongo database is ready. Note for easy use `docker`.

4. Start the application:

   ```bash
   mvn clean spring-boot:run
   ```

5. Run via docker

   ```bash
   docker compose up -f docker-compose-prod.yml up -d --build
   ```

The server will start at `http://localhost:8080`. Open the link in your browser to use the application.

---

## **Usage**

Once the application is running, you can:

- **Create invoices**: Add new invoices with relevant details.
- **Manage invoices**: View, update, or delete existing invoices.
- **Analyze data**: Gain insights from charts and invoice statistics.

---

## **Technologies Used**

- **Java Spring Boot**: Backend framework for building RESTful APIs and web applications.
- **Thymeleaf**: Templating engine for rendering dynamic HTML in Java web applications.
- **MongoDB**: NoSQL database for data storage.
- **Chart.js**: Interactive charting library for data visualization.
- **Tailwind CSS** and **DaisyUI**: For modern, responsive UI design.
- **Docker** and **Docker Compose**: For containerization and simplified deployment.

---

## **Contributing**

Contributions are welcome! Follow these steps to contribute:

1. Fork this repository.
2. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add your feature"
   ```
4. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Open a pull request to this repository.

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


