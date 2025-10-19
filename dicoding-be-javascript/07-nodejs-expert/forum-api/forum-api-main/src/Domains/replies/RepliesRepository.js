class RepliesRepository {
    async addReplyToComment(commentId, newReply, owner) {
        throw new Error('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async verifyReplyAccess(replyId, userId) {
        throw new Error('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async verifyReplyIsExist(replyId) {
        throw new Error('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async deleteReply(replyId) {
        throw new Error('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async repliesFromComment(commentId) {
        throw new Error('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }
}

module.exports = RepliesRepository;
