const ThreadRepository = require('../../../../Domains/threads/ThreadRepository');
const CommentRepository = require('../../../../Domains/comments/CommentRepository');
const RepliesRepository = require('../../../../Domains/replies/RepliesRepository');
const NewReply = require('../../../../Domains/replies/entities/NewReply');
const AddedReply = require('../../../../Domains/replies/entities/AddedReply');
const AddReplyUseCase = require('../AddReplyUseCase');

describe('AddReplyUseCase', () => {
    it('should orchestracting the add reply action correctly', async () => {
        // Arrange
        const payload = {
            content: 'sebuah balasan',
        };

        const owner = 'user-123';

        const mockAddedReply = new AddedReply({
            id: 'reply-123',
            content: payload.content,
            owner,
        });

        /** creating dependency of use case */
        const mockRepliesRepository = new RepliesRepository();
        const mockCommentRepository = new CommentRepository();
        const mockThreadRepository = new ThreadRepository();

        /** mocking needed function */
        mockRepliesRepository.addReplyToComment = jest.fn()
            .mockImplementation(() => Promise.resolve(mockAddedReply));
        mockCommentRepository.verifyCommentIsExist = jest.fn()
            .mockImplementation(() => Promise.resolve());
        mockThreadRepository.verifyThread = jest.fn()
            .mockImplementation(() => Promise.resolve());

        /** creating use case instance */
        const addCommentToThreadUseCase = new AddReplyUseCase({
            repliesRepository: mockRepliesRepository,
            commentRepository: mockCommentRepository,
            threadRepository: mockThreadRepository,
        });

        // Action
        const addedReply = await addCommentToThreadUseCase.execute('thread-123', 'comment-123', payload, owner);

        // Assert
        expect(addedReply).toStrictEqual(new AddedReply({
            id: 'reply-123', content: payload.content, owner,
        }));

        expect(mockThreadRepository.verifyThread).toBeCalledWith('thread-123');
        expect(mockCommentRepository.verifyCommentIsExist).toBeCalledWith('comment-123', 'thread-123');
        expect(mockRepliesRepository.addReplyToComment).toBeCalledWith('comment-123', new NewReply(payload), owner);
    });
});
