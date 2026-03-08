# **Finly - Invoicing Web Application**

**Finly** is a modern invoicing application designed to help businesses easily generate and manage invoices online. It includes features like invoice creation, management, and data visualization through interactive charts.

This project is part of the _Node.js Web Development For Beginners_ book, built using cutting-edge technologies like **Node.js**, **Express**, **MongoDB**, **EJS**, **Chart.js**, and styled with **Tailwind CSS** and **DaisyUI**.

---

## **Features**

- Create and manage invoices online.
- View detailed invoice statistics.
- Interactive data visualization with charts.
- User-friendly and responsive design powered by Tailwind CSS and DaisyUI.

---

## **Requirements**

- **Node.js** (v22 or higher)
- **MongoDB** (v7 or higher)

---

## **Getting Started**

1. Install dependencies:

   ```bash
   pnpm install
   ```

2. Copy the environment variables file:

   ```bash
   cp .env.sample .env
   ```

3. Ensure mongo database is ready. Note for easy use `docker`.

4. Start the application:

   ```bash
   npm run dev   # For development
   npm run start # For production
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

- **Node.js**: Backend runtime environment.
- **Express**: Web application framework for APIs and routing.
- **MongoDB**: NoSQL database for data storage.
- **Mongoose**: Object Data Modeling (ODM) library for MongoDB.
- **EJS**: Templating engine for rendering dynamic HTML.
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
