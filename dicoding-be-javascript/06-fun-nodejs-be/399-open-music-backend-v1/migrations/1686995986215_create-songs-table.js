/* eslint-disable camelcase */

exports.up = (pgm) => {
    pgm.createTable('songs', {
        id: {
            type: 'VARCHAR(50)',
            primaryKey: true,
        },
        title: {
            type: 'TEXT',
            notNull: true,
        },
        year: {
            type: 'INTEGER',
            notNull: true,
        },
        performer: {
            type: 'TEXT',
            notNull: true,
        },
        genre: {
            type: 'TEXT',
            notNull: true,
        },
        duration: {
            type: 'INTEGER',
        },
        album_id: {
            type: 'VARCHAR(50)',
        },
        created_at: {
            type: 'TIMESTAMP',
            notNull: true,
            default: pgm.func('CURRENT_TIMESTAMP'),
        },
        updated_at: {
            type: 'TIMESTAMP',
            notNull: true,
            default: pgm.func('CURRENT_TIMESTAMP'),
        },
    });

    pgm.addConstraint(
        'songs',
        'fk_songs.album_id',
        'FOREIGN KEY(album_id) REFERENCES albums(id) ON DELETE CASCADE',
    );
};

exports.down = (pgm) => {
    pgm.dropConstraint('songs', 'fk_songs.album_id');
    pgm.dropTable('songs');
};