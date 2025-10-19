/* eslint-disable camelcase */

exports.up = (pgm) => {
    pgm.createTable('albums', {
        id: {
            type: 'varchar(30)',
            primaryKey: true,
        },
        name: {
            type: 'varchar(50)',
            notNull: true,
        },
        year: {
            type: 'integer',
            notNull: true,
        },
        created_at: {
            type: 'text',
            notNull: true,
        },
        updated_at: {
            type: 'text',
            notNull: true,
        },
    });
};

exports.down = (pgm) => {
    pgm.dropTable('albums');
};
