const CommentRepository = require('../CommentRepository');

describe('CommentRepository', () => {
    it('should throw error when invoke unimplemented method', async () => {
        // Arrange
        const commentRepository = new CommentRepository();

        // Action & Assert
        await expect(commentRepository.addCommentToThread('', { content: '' }, ''))
            .rejects.toThrowError('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(commentRepository.verifyCommentAccess('', ''))
            .rejects.toThrowError('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(commentRepository.verifyCommentIsExist('', ''))
            .rejects.toThrowError('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(commentRepository.deleteCommentById(''))
            .rejects.toThrowError('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(commentRepository.commentsFromThread(''))
            .rejects.toThrowError('COMMENT_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    });
});
