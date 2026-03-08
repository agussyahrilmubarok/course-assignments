require('dotenv').config();

const bcrypt = require('bcrypt');

const log = require('./libs/logger');

const User = require('./models/user-model');
const Customer = require('./models/customer-model');
const Invoice = require('./models/invoice-model');

const { mongooseConnect } = require('./libs/mongodb');

const seedDatabase = async () => {
  try {
    await mongooseConnect();
    log.info('🔄 Seeding database started...');

    const existingUsers = await User.countDocuments();
    if (existingUsers > 0) {
      log.info('⚠️ Users already seeded. Skipping user seeding.');
    } else {
      const hashedPasswords = await Promise.all([
        bcrypt.hash('securepassword1', 10),
        bcrypt.hash('securepassword2', 10),
      ]);

      const users = await User.insertMany([
        { email: 'admin@example.com', password: hashedPasswords[0] },
        { email: 'user@example.com', password: hashedPasswords[1] },
      ]);
      log.info('✅ Users seeded successfully.', { users });
    }

    const existingCustomers = await Customer.countDocuments();
    if (existingCustomers > 0) {
      log.info('⚠️ Customers already seeded. Skipping customer seeding.');
    } else {
      const admin = await User.findOne({ email: 'admin@example.com' });
      const user = await User.findOne({ email: 'user@example.com' });

      const customers = await Customer.insertMany([
        {
          name: 'John Doe',
          email: 'johndoe@example.com',
          phone: '123456789',
          address: '123 Main St',
          owner: admin._id,
        },
        {
          name: 'Jane Smith',
          email: 'janesmith@example.com',
          phone: '987654321',
          address: '456 Elm St',
          owner: user._id,
        },
      ]);

      log.info('✅ Customers seeded successfully.', { customers });
    }

    const existingInvoices = await Invoice.countDocuments();
    if (existingInvoices > 0) {
      log.info('⚠️ Invoices already seeded. Skipping invoice seeding.');
    } else {
      const johnDoe = await Customer.findOne({ email: 'johndoe@example.com' });
      const janeSmith = await Customer.findOne({
        email: 'janesmith@example.com',
      });

      const admin = await User.findOne({ email: 'admin@example.com' });
      const user = await User.findOne({ email: 'user@example.com' });

      const invoices = await Invoice.insertMany([
        {
          amount: 100.0,
          date: '2024-12-19',
          status: 'Paid',
          owner: admin._id,
          customer: johnDoe._id,
        },
        {
          amount: 200.0,
          date: '2024-12-20',
          status: 'Unpaid',
          owner: user._id,
          customer: janeSmith._id,
        },
      ]);

      log.info('✅ Invoices seeded successfully.', { invoices });
    }

    log.info('✅ Seeding database completed successfully.');
    process.exit(0); // Exit script successfully
  } catch (err) {
    log.error('❌ Error seeding database.', {
      message: err.message,
      stack: err.stack,
    });
    process.exit(1); // Exit script with error
  }
};

seedDatabase();
