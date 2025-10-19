const RepliesRepository = require('../RepliesRepository');

describe('ThreadCommentRepliesRepository', () => {
    it('should throw error when invoke unimplemented method', async () => {
        // Arrange
        const repliesRepository = new RepliesRepository();

        // Action & Assert
        await expect(repliesRepository.addReplyToComment('', { content: '' }, ''))
            .rejects.toThrowError('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(repliesRepository.verifyReplyAccess('', ''))
            .rejects.toThrowError('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(repliesRepository.verifyReplyIsExist(''))
            .rejects.toThrowError('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(repliesRepository.deleteReply(''))
            .rejects.toThrowError('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(repliesRepository.repliesFromComment(''))
            .rejects.toThrowError('REPLIES_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    });
});
