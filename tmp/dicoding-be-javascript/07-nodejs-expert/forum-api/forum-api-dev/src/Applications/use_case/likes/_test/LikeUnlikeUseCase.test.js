const ThreadRepository = require('../../../../Domains/threads/ThreadRepository');
const CommentRepository = require('../../../../Domains/comments/CommentRepository');
const LikeRepository = require('../../../../Domains/likes/LikeRepository');
const LikeUnlikeUseCase = require('../LikeUnlikeUseCase');

describe('LikeUnlikeUseCase', () => {
    it('should orchestrating the unlike comment action correctly', async () => {
        // Arrange
        const useCasePayload = {
            threadId: 'thread-123',
            commentId: 'comment-123',
            userId: 'user-123',
        };

        /** creating dependency of use case */
        const mockThreadRepository = new ThreadRepository();
        const mockCommentRepository = new CommentRepository();
        const mockLikeRepository = new LikeRepository();

        /** mocking needed function */
        mockThreadRepository.verifyThread = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockCommentRepository.verifyCommentIsExist = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockLikeRepository.verifyLikeIsExists = jest.fn()
            .mockImplementation(() => Promise.resolve(1));
        mockLikeRepository.unLikeComment = jest.fn()
            .mockImplementation(() => Promise.resolve(1));

        /** creating use case instance */
        const likeUnlikeUseCase = new LikeUnlikeUseCase({
            threadRepository: mockThreadRepository,
            commentRepository: mockCommentRepository,
            likeRepository: mockLikeRepository,
        });

        // Action
        const isLike = await likeUnlikeUseCase.execute(useCasePayload);

        // Assert
        expect(isLike).toEqual(1);
        expect(mockThreadRepository.verifyThread).toBeCalledWith(useCasePayload.threadId);
        expect(mockCommentRepository.verifyCommentIsExist).toBeCalledWith(
            useCasePayload.commentId,
            useCasePayload.threadId,
        );
        expect(mockLikeRepository.verifyLikeIsExists).toBeCalledWith(useCasePayload);
        expect(mockLikeRepository.unLikeComment).toBeCalledWith(useCasePayload);
    });

    it('should orchestrating the like comment action correctly', async () => {
        // Arrange
        const useCasePayload = {
            threadId: 'thread-123',
            commentId: 'comment-123',
            userId: 'user-123',
        };

        const mockThreadRepository = new ThreadRepository();
        const mockCommentRepository = new CommentRepository();
        const mockLikeRepository = new LikeRepository();

        mockThreadRepository.verifyThread = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockCommentRepository.verifyCommentIsExist = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockLikeRepository.verifyLikeIsExists = jest.fn()
            .mockImplementation(() => Promise.resolve(0));
        mockLikeRepository.addLikeToComment = jest.fn()
            .mockImplementation(() => Promise.resolve(1));

        const likeUnlikeUseCase = new LikeUnlikeUseCase({
            threadRepository: mockThreadRepository,
            commentRepository: mockCommentRepository,
            likeRepository: mockLikeRepository,
        });

        // Action
        const isUnlike = await likeUnlikeUseCase.execute(useCasePayload);

        // Assert
        expect(isUnlike).toEqual(1);
        expect(mockThreadRepository.verifyThread)
            .toBeCalledWith(useCasePayload.threadId);
        expect(mockCommentRepository.verifyCommentIsExist)
            .toBeCalledWith(useCasePayload.commentId, useCasePayload.threadId);
        expect(mockLikeRepository.verifyLikeIsExists).toBeCalledWith(useCasePayload);
        expect(mockLikeRepository.addLikeToComment).toBeCalledWith(useCasePayload);
    });
});
