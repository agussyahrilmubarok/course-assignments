exports.up = (pgm) => {
  pgm.createTable('festivals', {
    id: {
      type: 'VARCHAR(50)',
      primaryKey: true,
    },
    title: {
      type: 'TEXT',
      notNull: true,
    },
    description: {
      type: 'TEXT',
      notNull: true,
    },
    facilities: {
      type: 'TEXT[]',
      notNull: true,
    },
    refundable: {
      type: 'BOOLEAN',
      notNull: true,
    },
    banner: {
      type: 'TEXT',
      notNull: true,
    },
    price: {
      type: 'DECIMAL(10,2)',
      notNull: true,
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('festivals');
};
