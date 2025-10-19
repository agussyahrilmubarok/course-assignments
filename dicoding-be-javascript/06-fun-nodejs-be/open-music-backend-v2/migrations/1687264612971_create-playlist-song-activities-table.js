/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('playlist_song_activities', {
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
            notNull: true,
        },
        user_id: {
            type: 'varchar(30)',
            notNull: true,
        },
        action: {
            type: 'text',
            notNull: true,
        },
        time: {
            type: 'text',
            notNull: true,
        },
    });
};

exports.down = (pgm) => {
    pgm.dropTable('playlist_song_activities');
};
