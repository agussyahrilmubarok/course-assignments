/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('playlists', {
        id: {
            type: 'varchar(30)',
            primaryKey: true,
        },
        name: {
            type: 'varchar(50)',
            notNull: true,
        },
        owner: {
            type: 'varchar(30)',
            references: '"users"',
            onDelete: 'cascade',
            onUpdate: 'cascade',
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
    pgm.dropTable('playlists');
};
