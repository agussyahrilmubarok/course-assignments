/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('comment_likes', {
        id: {
            type: 'VARCHAR(50)',
            primaryKey: true,
        },
        comment_id: {
            type: 'VARCHAR(50)',
            notNull: true,
        },
        owner: {
            type: 'VARCHAR(50)',
            notNull: true,
        },
        date: {
            type: 'TIMESTAMP',
            notNull: true,
            default: pgm.func('current_timestamp'),
        },
    });

    // add constraint unique
    pgm.addConstraint(
        'comment_likes',
        'no_duplicate_likes',
        'UNIQUE(comment_id, owner)',
    );

    // add constraint foreign key to comment_id references to comment.id
    pgm.addConstraint(
        'comment_likes',
        'fk_comment_likes.comment_id_comments.id',
        'FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE',
    );

    // add constraint foreign key to owner references to users.id
    pgm.addConstraint(
        'comment_likes',
        'fk_comment_likes.owner_users.id',
        'FOREIGN KEY(owner) REFERENCES users(id) ON DELETE CASCADE',
    );
};

exports.down = (pgm) => {
    pgm.dropConstraint('comment_likes', 'fk_comment_likes.comment_id_comments.id');
    pgm.dropConstraint('comment_likes', 'fk_comment_likes.owner_users.id');
    pgm.dropConstraint('comment_likes', 'no_duplicate_likes');
    pgm.dropTable('comment_likes');
};
