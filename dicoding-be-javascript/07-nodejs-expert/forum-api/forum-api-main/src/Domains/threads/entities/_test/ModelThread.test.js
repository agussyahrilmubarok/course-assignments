const ModelThread = require('../ModelThread');

describe('a Thread entities', () => {
    it('should throw error when payload did not contain needed property', () => {
        // Arrange
        const payload = {
            title: 'sebuah thread',
        };

        // Action and Assert
        expect(() => new ModelThread(payload)).toThrowError('THREAD.NOT_CONTAIN_NEEDED_PROPERTY');
    });

    it('should throw error when payload did not meet data type specification', () => {
        // Arrange
        const payload = {
            id: 'thread-123',
            title: 'sebuah thread',
            body: 'sebuah body thread',
            date: Date.now(),
            username: 'jhon doe',
        };

        // Action and Assert
        expect(() => new ModelThread(payload)).toThrowError('THREAD.NOT_MEET_DATA_TYPE_SPECIFICATION');
    });

    it('should create Thread object correctly', () => {
        // Arrange
        const payload = {
            id: 'thread-123',
            title: 'sebuah thread',
            body: 'sebuah body thread',
            date: new Date(),
            username: 'jhon doe',
        };

        // Action
        const modelThread = new ModelThread(payload);

        // Assert
        expect(modelThread).toBeInstanceOf(ModelThread);
        expect(modelThread.id).toStrictEqual(payload.id);
        expect(modelThread.title).toStrictEqual(payload.title);
        expect(modelThread.body).toStrictEqual(payload.body);
        expect(modelThread.date).toStrictEqual(payload.date);
        expect(modelThread.username).toStrictEqual(payload.username);
    });
});
