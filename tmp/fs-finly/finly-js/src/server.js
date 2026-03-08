require('dotenv').config();

const express = require('express');
const session = require('express-session');
const flash = require('connect-flash');
const morgan = require('morgan');
const log = require('./libs/logger');
const authRouter = require('./routes/auth-router');
const dashboardRouter = require('./routes/dashboard-router');

const { mongooseConnect } = require('./libs/mongodb');
const { config } = require('./libs/config');
const { verifyUser } = require('./libs/middleware');

const app = express();

app.set('views', './src/views');
app.set('view engine', 'ejs');

app.use(
  morgan('combined', { stream: { write: (msg) => log.info(msg.trim()) } })
);
app.use(express.static('./public'));
app.use(express.urlencoded({ extended: false }));
app.use(flash());
app.use(
  session({
    secret: config.AUTH_SECRET,
    saveUninitialized: true,
    resave: false,
  })
);

app.use('/', authRouter);
app.use('/dashboard', verifyUser, dashboardRouter);

app.use((req, res) => {
  log.warn('404 Not Found', { path: req.originalUrl });
  res.status(404).render('index', {
    title: 'Not Found',
    message: 'Not Found',
  });
});

app.use((err, req, res, next) => {
  log.error('Unhandled error occurred:', {
    message: err.message,
    stack: err.stack,
    action: 'Unhandled error in request',
  });
  res.status(500).send('Internal Server Error');
});

const port = config.APP_PORT;

mongooseConnect()
  .then(() => {
    app.listen(port, () => {
      log.info(`✅ Server is running on port ${port}`, { port });
    });
  })
  .catch((err) => {
    log.error('❌ Error connecting to the database.', {
      message: err.message,
      stack: err.stack,
    });
    process.exit(1);
  });
