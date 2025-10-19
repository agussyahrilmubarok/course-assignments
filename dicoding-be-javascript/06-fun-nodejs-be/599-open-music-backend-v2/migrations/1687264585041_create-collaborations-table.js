/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('collaborations', {
        id: {
            type: 'varchar(30)',
            primaryKey: true,
        },
        playlist_id: {
            type: 'varchar(30)',
            references: '"playlists"',
            onDelete: 'cascade',
            onUpdate: 'cascade',
            unique: true,
        },
        user_id: {
            type: 'varchar(30)',
            references: '"users"',
            onDelete: 'cascade',
            onUpdate: 'cascade',
            unique: true,
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
    pgm.dropTable('collaborations');
};
