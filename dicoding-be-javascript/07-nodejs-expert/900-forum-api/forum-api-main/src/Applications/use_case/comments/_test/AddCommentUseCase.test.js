const ThreadRepository = require('../../../../Domains/threads/ThreadRepository');
const CommentRepository = require('../../../../Domains/comments/CommentRepository');
const NewComment = require('../../../../Domains/comments/entities/NewComment');
const AddedComment = require('../../../../Domains/comments/entities/AddedComment');
const AddCommentUseCase = require('../AddCommentUseCase');

describe('AddCommentUseCase', () => {
    it('should orchestracting the add comment action correctly', async () => {
        // Arrange
        const payload = {
            content: 'sebuah comment',
        };

        const owner = 'user-123';

        const mockAddedComment = new AddedComment({
            id: 'comment-123',
            content: payload.content,
            owner,
        });

        /** creating dependency of use case */
        const mockCommentRepository = new CommentRepository();
        const mockThreadRepository = new ThreadRepository();

        /** mocking needed function */
        mockCommentRepository.addCommentToThread = jest.fn()
            .mockImplementation(() => Promise.resolve(mockAddedComment));
        mockThreadRepository.verifyThread = jest.fn()
            .mockImplementation(() => Promise.resolve());

        /** creating use case instance */
        const addCommentUseCase = new AddCommentUseCase({
            commentRepository: mockCommentRepository,
            threadRepository: mockThreadRepository,
        });

        // Action
        const addedComment = await addCommentUseCase.execute('thread-123', payload, owner);

        // Assert
        expect(addedComment).toStrictEqual(new AddedComment({
            id: 'comment-123', content: 'sebuah comment', owner,
        }));

        expect(mockThreadRepository.verifyThread).toBeCalledWith('thread-123');
        expect(mockCommentRepository.addCommentToThread).toBeCalledWith('thread-123', new NewComment(payload), owner);
    });
});
