class CommentRepository {
    async addCommentToThread(threadId, newComment, owner) {
        throw new Error('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async verifyCommentAccess(commentId, userId) {
        throw new Error('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async verifyCommentIsExist(commentId, threadId) {
        throw new Error('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async deleteCommentById(commentId) {
        throw new Error('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async commentsFromThread(threadId) {
        throw new Error('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }
}

module.exports = CommentRepository;
