const CommentRepository = require('../../../../Domains/comments/CommentRepository');
const DeleteCommentUseCase = require('../DeleteCommentUseCase');

describe('DeleteCommentUseCase', () => {
    it('should orchestracting the delete comment action correctly', async () => {
        // Arrange
        const useCasePayload = {
            threadId: 'thread-123',
            commentId: 'comment-123',
            userId: 'user-123',
        };

        /** creating dependency of use case */
        const mockCommentRepository = new CommentRepository();

        /** mocking needed function */
        mockCommentRepository.deleteCommentById = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockCommentRepository.verifyCommentAccess = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockCommentRepository.verifyCommentIsExist = jest.fn()
            .mockImplementation(() => Promise.resolve());

        /** creating use case instance */
        const deleteCommentUseCase = new DeleteCommentUseCase({
            commentRepository: mockCommentRepository,
        });

        // Action
        await deleteCommentUseCase.execute({
            threadId: useCasePayload.threadId,
            commentId: useCasePayload.commentId,
            userId: useCasePayload.userId,
        });

        // Assert
        expect(mockCommentRepository.verifyCommentIsExist).toBeCalledWith(
            useCasePayload.commentId,
            useCasePayload.threadId,
        );
        expect(mockCommentRepository.verifyCommentAccess).toBeCalledWith(
            useCasePayload.commentId,
            useCasePayload.userId,
        );
        expect(mockCommentRepository.deleteCommentById).toBeCalledWith(
            useCasePayload.commentId,
        );
    });
});
