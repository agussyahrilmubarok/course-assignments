/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("comments", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    thread_id: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    content: {
      type: "TEXT",
      notNull: true,
    },
    owner: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    is_deleted: {
      type: "BOOLEAN",
      notNull: true,
      default: false,
    },
    date: {
      type: "TIMESTAMP",
      notNull: true,
      default: pgm.func("current_timestamp"),
    },
  });

  pgm.addConstraint(
    "comments",
    "fk_comments.threads.thread_id",
    "FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE"
  );
  pgm.addConstraint(
    "comments",
    "fk_comments.users.owner",
    "FOREIGN KEY(owner) REFERENCES users(id) ON DELETE CASCADE"
  );
  pgm.createIndex("comments", "owner");
  pgm.createIndex("comments", "thread_id");
};

exports.down = (pgm) => {
  pgm.dropConstraint("comments", "fk_comments.threads.thread_id");
  pgm.dropConstraint("comments", "fk_comments.users.owner");
  pgm.dropTable("comments");
};
