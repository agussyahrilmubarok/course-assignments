const ThreadRepository = require('../ThreadRepository');

describe('ThreadRepository', () => {
    it('should throw error when invoke unimplemented method', async () => {
        // Arrange
        const threadRepository = new ThreadRepository();

        // Action & Assert
        await expect(threadRepository.addThread({ title: '', body: '' }, ''))
            .rejects.toThrowError('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(threadRepository.getThreadById(''))
            .rejects.toThrowError('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');

        await expect(threadRepository.verifyThread(''))
            .rejects.toThrowError('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    });
});
