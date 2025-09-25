const log = require('../libs/logger');
const Customer = require('../models/customer-model');
const Invoice = require('../models/invoice-model');

const { CurrencyUSDollar } = require('../libs/currency');

const showDashboard = async (req, res) => {
  const query = { owner: req.session.userId };

  try {
    log.info('Loading dashboard', { userId: req.session.userId });

    const invoiceCount = await Invoice.countDocuments(query);
    const customerCount = await Customer.countDocuments(query);
    const allInvoices = await Invoice.find(query).populate({
      path: 'customer',
      model: Customer,
      select: '_id name',
    });

    // Calculate total paid and total pending
    const totalPaid = allInvoices.reduce((sum, invoice) => {
      return invoice.status === 'paid' ? sum + invoice.amount : sum;
    }, 0);

    const totalPending = allInvoices.reduce((sum, invoice) => {
      return invoice.status === 'pending' ? sum + invoice.amount : sum;
    }, 0);

    // Sort invoices by date (newest first)
    allInvoices.sort((a, b) => new Date(b.date) - new Date(a.date));
    const latestInvoices = allInvoices.slice(0, 5);

    // Prepare revenue data for the last 6 months
    const revenueData = [];
    for (let i = 0; i < 6; i++) {
      const today = new Date();
      today.setMonth(today.getMonth() - i);
      const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
      const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0);
      const month = today.toLocaleString('default', { month: 'short' });

      const revenueForMonth = allInvoices
        .filter((invoice) => {
          return (
            new Date(invoice.date) >= firstDay &&
            new Date(invoice.date) <= lastDay
          );
        })
        .reduce((total, invoice) => total + invoice.amount, 0);

      revenueData.unshift({ month, revenue: revenueForMonth });
    }

    log.info('Dashboard data retrieved successfully', {
      invoiceCount,
      customerCount,
      totalPaid,
      totalPending,
      latestInvoicesCount: latestInvoices.length,
    });

    res.render('pages/dashboard', {
      title: 'Dashboard',
      latestInvoices,
      revenueData: JSON.stringify(revenueData),
      invoiceCount,
      customerCount,
      totalPaid,
      totalPending,
      CurrencyUSDollar,
      info: req.flash('info')[0],
    });
  } catch (err) {
    log.error('Error loading dashboard', {
      error: err.message,
      stack: err.stack,
      userId: req.session.userId,
    });
    res.status(500).send('Error loading dashboard');
  }
};

module.exports = { showDashboard };
