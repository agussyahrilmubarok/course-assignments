const NewThread = require('../NewThread');

describe('a NewThread entities', () => {
    it('should throw error when payload did not contain needed property', () => {
        // Arrange
        const payload = {
            title: 'sebuah thread',
        };

        // Action and Assert
        expect(() => new NewThread(payload)).toThrowError('NEW_THREAD.NOT_CONTAIN_NEEDED_PROPERTY');
    });

    it('should throw error when payload did not meet data type specification', () => {
        // Arrange
        const payload = {
            title: 'sebuah thread',
            body: 123,
        };

        // Action and Assert
        expect(() => new NewThread(payload)).toThrowError('NEW_THREAD.NOT_MEET_DATA_TYPE_SPECIFICATION');
    });

    it('should create Thread object correctly', () => {
        // Arrange
        const payload = {
            title: 'sebuah thread',
            body: 'sebuah body thread',
        };

        // Action
        const newThread = new NewThread(payload);

        // Assert
        expect(newThread).toBeInstanceOf(NewThread);
        expect(newThread.title).toStrictEqual(payload.title);
        expect(newThread.body).toStrictEqual(payload.body);
    });
});
