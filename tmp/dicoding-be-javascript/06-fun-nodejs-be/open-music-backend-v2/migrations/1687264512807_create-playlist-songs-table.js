/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('playlist_songs', {
        id: {
            type: 'varchar(30)',
            primaryKey: true,
        },
        playlist_id: {
            type: 'varchar(30)',
            references: '"playlists"',
            onDelete: 'cascade',
            onUpdate: 'cascade',
        },
        song_id: {
            type: 'varchar(30)',
            references: '"songs"',
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
    pgm.dropTable('playlist_songs');
};

