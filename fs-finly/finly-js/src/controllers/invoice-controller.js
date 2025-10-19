const log = require('../libs/logger');
const Customer = require('../models/customer-model');
const Invoice = require('../models/invoice-model');

const { body, validationResult } = require('express-validator');
const { CurrencyUSDollar } = require('../libs/currency');

const validateInvoice = [
  body('customer', 'Select the Customer').notEmpty(),
  body('amount', 'Amount must not be empty').notEmpty(),
  body('date', 'Due Date must not be empty').notEmpty(),
  body('status', 'Select the Status').notEmpty(),
];

const populateInvoices = async (query, search) => {
  const populateOptions = {
    path: 'customer',
    model: Customer,
    select: '_id name',
  };
  if (search) {
    populateOptions['match'] = { name: { $regex: search, $options: 'i' } };
  }

  const invoices = await query.populate(populateOptions);
  if (!Array.isArray(invoices)) {
    return [];
  }

  return invoices.filter((invoice) => invoice.customer != null);
};

const showInvoices = async (req, res) => {
  const query = { owner: req.session.userId };
  const { search } = req.query;

  try {
    log.info('Fetching invoices', { userId: req.session.userId, search });

    const invoices = await populateInvoices(Invoice.find(query), search);

    log.info('Invoices fetched successfully', {
      invoiceCount: invoices.length,
    });

    res.render('pages/invoices', {
      title: 'Invoices',
      type: 'data',
      invoices: invoices || [], // Pastikan invoices memiliki nilai default
      CurrencyUSDollar,
      info: req.flash('info')[0],
    });
  } catch (err) {
    log.error('Error fetching invoices', {
      error: err.message,
      stack: err.stack,
    });
    res.status(500).send('Error fetching invoices');
  }
};

const getCustomers = async (req, res, next) => {
  const customersQuery = { owner: req.session.userId };
  try {
    const customers = await Customer.find(customersQuery);
    req.customers = customers;
    next();
  } catch (err) {
    log.error('Error fetching customers for invoice', {
      error: err.message,
      stack: err.stack,
    });
    next(err);
  }
};

const editInvoice = async (req, res) => {
  const invoiceId = req.params.id;

  try {
    log.info('Fetching invoice data for editing', {
      invoiceId,
      userId: req.session.userId,
    });

    const invoice = await Invoice.findById(invoiceId).populate({
      path: 'customer',
      model: Customer,
      select: '_id name',
    });

    const { customers } = req;

    if (!invoice) {
      log.warn('Invoice not found for editing', { invoiceId });
      return res.status(404).send('Invoice not found');
    }

    log.info('Invoice data fetched for editing', { invoiceId });

    res.render('pages/invoices', {
      title: 'Edit Invoice',
      type: 'form',
      formAction: 'edit',
      customers,
      invoice: req.flash('data')[0] || invoice,
      errors: req.flash('errors'),
    });
  } catch (err) {
    log.error('Error editing invoice', {
      error: err.message,
      stack: err.stack,
      invoiceId,
    });
    res.status(500).send('Error editing invoice');
  }
};

const createInvoice = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Invoice creation validation failed', {
      errors,
      userId: req.session.userId,
    });
    return res.redirect('create');
  }

  const newInvoice = req.body;
  newInvoice.owner = req.session.userId;

  try {
    const createdInvoice = await Invoice.create(newInvoice);

    log.info('Invoice created successfully', {
      invoiceId: createdInvoice._id,
      userId: req.session.userId,
    });

    req.flash('info', { message: 'New Invoice Created', type: 'success' });
    res.redirect('/dashboard/invoices');
  } catch (err) {
    log.error('Error creating invoice', {
      error: err.message,
      stack: err.stack,
      userId: req.session.userId,
    });
    req.flash('info', {
      message: 'An error occurred while creating the invoice',
      type: 'error',
    });
    res.redirect('create');
  }
};

const updateInvoice = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Invoice update validation failed', {
      errors,
      userId: req.session.userId,
    });
    return res.redirect('edit');
  }

  const invoiceId = req.params.id;
  const data = req.body;

  try {
    const updatedInvoice = await Invoice.findByIdAndUpdate(invoiceId, data);

    if (!updatedInvoice) {
      log.warn('Invoice not found for update', { invoiceId });
      return res.status(404).send('Invoice not found');
    }

    log.info('Invoice updated successfully', {
      invoiceId,
      userId: req.session.userId,
    });

    req.flash('info', { message: 'Invoice Updated', type: 'success' });
    res.redirect('/dashboard/invoices');
  } catch (err) {
    log.error('Error updating invoice', {
      error: err.message,
      stack: err.stack,
      invoiceId,
    });
    req.flash('info', {
      message: 'An error occurred while updating the invoice',
      type: 'error',
    });
    res.redirect('edit');
  }
};

const deleteInvoice = async (req, res) => {
  const invoiceId = req.params.id;

  try {
    const deletedInvoice = await Invoice.findByIdAndDelete(invoiceId);
    if (!deletedInvoice) {
      log.warn('Invoice not found for deletion', { invoiceId });
      return res.status(404).send('Invoice not found');
    }

    log.info('Invoice deleted successfully', {
      invoiceId,
      userId: req.session.userId,
    });

    req.flash('info', { message: 'Invoice Deleted', type: 'success' });
    res.redirect('/dashboard/invoices');
  } catch (err) {
    log.error('Error deleting invoice', {
      error: err.message,
      stack: err.stack,
      invoiceId,
    });
    req.flash('info', {
      message: 'An error occurred while deleting the invoice',
      type: 'error',
    });
    res.redirect('/dashboard/invoices');
  }
};

module.exports = {
  validateInvoice,
  showInvoices,
  editInvoice,
  deleteInvoice,
  updateInvoice,
  createInvoice,
  getCustomers,
};
