/* istanbul ignore file */
const pool = require('../src/Infrastructures/database/postgres/pool');

const CommentLikeTableTestHelper = {
    async addLikeToComment({
        id = 'like-123',
        commentId = 'comment-123',
        owner = 'user-123',
    }) {
        const query = {
            text: 'INSERT INTO comment_likes VALUES($1, $2, $3)',
            values: [id, commentId, owner],
        };

        await pool.query(query);
    },

    async findLikeCommentById(LikeId) {
        const query = {
            text: 'SELECT * FROM comment_likes WHERE id = $1',
            values: [LikeId],
        };

        const result = await pool.query(query);

        return result.rows;
    },

    async unLikeComment(LikeId) {
        const query = {
            text: 'DELETE comment_likes WHERE id = $1',
            values: [LikeId],
        };

        await pool.query(query);
    },

    async cleanTable() {
        await pool.query('DELETE FROM comment_likes WHERE 1=1');
    },
};

module.exports = CommentLikeTableTestHelper;