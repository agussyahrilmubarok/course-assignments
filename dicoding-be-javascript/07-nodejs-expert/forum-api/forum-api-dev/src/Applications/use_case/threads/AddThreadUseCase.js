const NewThread = require('../../../Domains/threads/entities/NewThread');

class AddThreadUseCase {
    constructor({ threadRepository }) {
        this._threadRepository = threadRepository;
    }

    async execute(payload, owner) {
        const thread = new NewThread(payload);

        return this._threadRepository.addThread(thread, owner);
    }
}

module.exports = AddThreadUseCase;
