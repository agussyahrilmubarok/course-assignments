const RepliesRepository = require('../../../../Domains/replies/RepliesRepository');
const DeleteReplyUseCase = require('../DeleteReplyUseCase');

describe('DeleteReplyUseCase', () => {
    it('should orchestracting the delete reply action correctly', async () => {
        // Arrange
        const useCaseDeletePayload = {
            threadId: 'thread-123',
            commentId: 'comment-123',
            replyId: 'reply-123',
            userId: 'user-123',
        };

        /** creating dependency of use case */
        const mockRepliesRepository = new RepliesRepository();

        /** mocking needed function */
        mockRepliesRepository.deleteReply = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockRepliesRepository.verifyReplyAccess = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockRepliesRepository.verifyReplyIsExist = jest.fn()
            .mockImplementation(() => Promise.resolve());

        /** creating use case instance */
        const deleteReplyUseCase = new DeleteReplyUseCase({
            repliesRepository: mockRepliesRepository,
        });

        // Action
        await deleteReplyUseCase.execute(useCaseDeletePayload);

        // Assert
        expect(mockRepliesRepository.verifyReplyIsExist).toBeCalledWith(
            useCaseDeletePayload.replyId,
        );
        expect(mockRepliesRepository.verifyReplyAccess).toBeCalledWith(
            useCaseDeletePayload.replyId,
            useCaseDeletePayload.userId,
        );
        expect(mockRepliesRepository.deleteReply).toBeCalledWith(useCaseDeletePayload.replyId);
    });
});
