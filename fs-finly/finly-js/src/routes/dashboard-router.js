const express = require('express');
const router = express.Router();

const customersRouter = require('./customer-router');
const invoicesRouter = require('./invoice-router');

const { showDashboard } = require('../controllers/dashboard-controller');

router.get('/', showDashboard);

router.use('/customers', customersRouter);

router.use('/invoices', invoicesRouter);

module.exports = router;
