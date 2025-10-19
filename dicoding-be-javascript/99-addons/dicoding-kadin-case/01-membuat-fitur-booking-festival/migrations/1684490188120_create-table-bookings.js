exports.up = (pgm) => {
  pgm.createTable('bookings', {
    id: {
      type: 'VARCHAR(50)',
      primaryKey: true,
    },
    user_id: {
      type: 'VARCHAR(50)',
      notNull: true,
    },
    festival_id: {
      type: 'VARCHAR(50)',
      notNull: true,
    },
    booking_date: {
      type: 'DATE',
      notNull: true,
    },
    quantity: {
      type: 'SMALLINT',
      notNull: true,
    },
    total_price: {
      type: 'DECIMAL(10,2)',
    },
    status: {
      type: 'SMALLINT',
      notNull: true,
    },
    created_at: {
      type: 'DATE',
      notNull: true,
    },
    confirmation_code: {
      type: 'VARCHAR(50)',
    },
  });

  pgm.addConstraint('bookings', 'fk_bookings.user_id_users.id', 'FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE');
  pgm.addConstraint('bookings', 'fk_bookings.festival_id_festivals.id', 'FOREIGN KEY(festival_id) REFERENCES festivals(id) ON DELETE CASCADE');
};

exports.down = (pgm) => {
  pgm.dropTable('bookings');
};
