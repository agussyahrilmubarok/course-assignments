const { Schema, model } = require('mongoose');

const InvoiceSchema = new Schema({
  amount: { type: Number, required: true },
  date: { type: String, required: true },
  status: { type: String, required: true },
  owner: { type: Schema.Types.ObjectId, ref: 'users' },
  customer: { type: Schema.Types.ObjectId, ref: 'customers' },
});

const Invoice = model('invoices', InvoiceSchema);

module.exports = Invoice;
