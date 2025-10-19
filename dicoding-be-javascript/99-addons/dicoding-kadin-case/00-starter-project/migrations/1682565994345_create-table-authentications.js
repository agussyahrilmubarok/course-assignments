exports.up = (pgm) => {
  pgm.createTable('authentications', {
    refresh_token: {
      type: 'TEXT',
      notNull: true,
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable('authentications');
};
