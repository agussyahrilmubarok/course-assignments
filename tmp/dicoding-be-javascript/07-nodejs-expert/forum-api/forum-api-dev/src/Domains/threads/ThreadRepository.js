class ThreadRepository {
    async addThread(newThread, owner) {
        throw new Error('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async getThreadById(threadId) {
        throw new Error('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }

    async verifyThread(threadId) {
        throw new Error('THREAD_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    }
}

module.exports = ThreadRepository;
