const log = require('../libs/logger');
const Customer = require('../models/customer-model');
const Invoice = require('../models/invoice-model');

const { body, validationResult } = require('express-validator');

const validateCustomer = [
  body('name', 'Name must not be empty').notEmpty(),
  body('email', 'Email must not be empty').notEmpty(),
  body('phone', 'Phone must not be empty').notEmpty(),
  body('address', 'Address must not be empty').notEmpty(),
];

const showCustomers = async (req, res) => {
  const query = { owner: req.session.userId };
  const { search } = req.query;
  if (search) {
    query['$or'] = [
      { name: { $regex: search, $options: 'i' } },
      { email: { $regex: search, $options: 'i' } },
      { phone: { $regex: search, $options: 'i' } },
      { address: { $regex: search, $options: 'i' } },
    ];
  }

  try {
    log.info('Fetching customer data', {
      userId: req.session.userId,
      search,
    });

    const customers = await Customer.find(query);

    log.info('Customers fetched successfully', {
      customerCount: customers.length,
    });

    res.render('pages/customers', {
      title: 'Customers',
      type: 'data',
      customers,
      info: req.flash('info')[0],
    });
  } catch (err) {
    log.error('Error loading customers data', {
      error: err.message,
      stack: err.stack,
    });
    res.status(500).send('Error loading customers data');
  }
};

const createCustomer = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Customer creation validation failed', {
      errors,
      userId: req.session.userId,
    });
    return res.redirect('create');
  }

  const newCustomer = req.body;
  newCustomer.owner = req.session.userId;

  try {
    await Customer.create(newCustomer);
    log.info('New customer created successfully', {
      customer: newCustomer,
      userId: req.session.userId,
    });

    req.flash('info', { message: 'Customer Created', type: 'success' });
    res.redirect('/dashboard/customers');
  } catch (err) {
    log.error('Error creating customer', {
      error: err.message,
      stack: err.stack,
      userId: req.session.userId,
    });
    res.status(500).send('Error creating customer');
  }
};

const editCustomer = async (req, res) => {
  const customerId = req.params.id;

  try {
    const customer = await Customer.findById(customerId);
    if (!customer) {
      log.warn('Customer not found for editing', {
        customerId,
        userId: req.session.userId,
      });
      return res.status(404).send('Customer not found');
    }

    log.info('Customer data fetched for editing', {
      customerId,
      userId: req.session.userId,
    });

    res.render('pages/customers', {
      title: 'Edit Customer',
      type: 'form',
      formAction: 'edit',
      customer: req.flash('data')[0] || customer,
      errors: req.flash('errors'),
    });
  } catch (err) {
    log.error('Error fetching customer data for editing', {
      error: err.message,
      stack: err.stack,
      customerId,
    });
    res.status(500).send('Error fetching customer data for editing');
  }
};

const updateCustomer = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Customer update validation failed', {
      errors,
      userId: req.session.userId,
    });
    return res.redirect('edit');
  }

  const customerId = req.params.id;
  const customerData = req.body;

  try {
    const updatedCustomer = await Customer.findByIdAndUpdate(
      customerId,
      customerData,
      { new: true }
    );
    if (!updatedCustomer) {
      log.warn('Customer not found for update', {
        customerId,
        userId: req.session.userId,
      });
      return res.status(404).send('Customer not found');
    }

    log.info('Customer updated successfully', {
      customerId,
      userId: req.session.userId,
      updatedCustomer,
    });

    req.flash('info', { message: 'Customer Updated', type: 'success' });
    res.redirect('/dashboard/customers');
  } catch (err) {
    log.error('Error updating customer', {
      error: err.message,
      stack: err.stack,
      customerId,
    });
    res.status(500).send('Error updating customer');
  }
};

const deleteCustomer = async (req, res) => {
  const customerId = req.params.id;

  try {
    const deletedInvoices = await Invoice.deleteMany({ customer: customerId });
    const deletedCustomer = await Customer.findByIdAndDelete(customerId);
    if (!deletedCustomer) {
      log.warn('Customer not found for deletion', {
        customerId,
        userId: req.session.userId,
      });
      return res.status(404).send('Customer not found');
    }

    log.info('Customer and associated invoices deleted successfully', {
      customerId,
      userId: req.session.userId,
    });

    req.flash('info', { message: 'Customer Deleted', type: 'success' });
    res.redirect('/dashboard/customers');
  } catch (err) {
    log.error('Error deleting customer', {
      error: err.message,
      stack: err.stack,
      customerId,
    });
    res.status(500).send('Error deleting customer');
  }
};

module.exports = {
  validateCustomer,
  showCustomers,
  createCustomer,
  editCustomer,
  updateCustomer,
  deleteCustomer,
};
