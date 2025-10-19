const ModelReply = require('../ModelReply');

describe('a Reply entities', () => {
    it('should throw error when payload did not contain needed property', () => {
        // AAA
        expect(() => new ModelReply({})).toThrowError('COMMENT.NOT_CONTAIN_NEEDED_PROPERTY');
    });

    it('should throw error when payload did not meet data type specification', () => {
        // Arrange
        const payload = {
            id: 'reply-123',
            username: 'dicoding',
            date: Date.now(),
            content: 'sebuah balasan',
            isDeleted: false,
        };

        // Action and Assert
        expect(() => new ModelReply(payload)).toThrowError('COMMENT.NOT_MEET_DATA_TYPE_SPECIFICATION');
    });

    it('should create Reply object correctly', () => {
        // Arrange
        const payload = {
            id: 'reply-123',
            username: 'dicoding',
            date: new Date(),
            content: 'sebuah balasan',
            isDeleted: false,
        };

        // Action
        const modelReply = new ModelReply(payload);

        // Assert
        expect(modelReply).toBeInstanceOf(ModelReply);
        expect(modelReply.id).toStrictEqual(payload.id);
        expect(modelReply.username).toStrictEqual(payload.username);
        expect(modelReply.date).toStrictEqual(payload.date);
        expect(modelReply.content).toStrictEqual(payload.content);
    });

    it('should create Reply object correctly with modified content when isDeleted argument is true', () => {
        // Arrange
        const payload = {
            id: 'reply-123',
            username: 'dicoding',
            date: new Date(),
            content: 'sebuah balasan',
            isDeleted: true,
        };

        // Action
        const modelReply = new ModelReply(payload);

        // Assert
        expect(modelReply).toBeInstanceOf(ModelReply);
        expect(modelReply.id).toStrictEqual(payload.id);
        expect(modelReply.username).toStrictEqual(payload.username);
        expect(modelReply.date).toStrictEqual(payload.date);
        expect(modelReply.content).toStrictEqual('**balasan telah dihapus**');
    });
});
