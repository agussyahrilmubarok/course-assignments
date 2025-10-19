const ModelComment = require('../ModelComment');

describe('a Comment entities', () => {
    it('should throw error when payload did not contain needed property', () => {
        // Arrange, Action & Assert
        expect(() => new ModelComment({})).toThrowError('COMMENT.NOT_CONTAIN_NEEDED_PROPERTY');
    });

    it('should throw error when payload did not meet data type specification', () => {
        // Arrange
        const payload = {
            id: 'comment-123',
            username: 'dicoding',
            date: Date.now(),
            content: 'sebuah comment',
            isDeleted: false,
        };

        // Action and Assert
        expect(() => new ModelComment(payload)).toThrowError('COMMENT.NOT_MEET_DATA_TYPE_SPECIFICATION');
    });

    it('should create Comment object correctly', () => {
        // Arrange
        const payload = {
            id: 'comment-123',
            username: 'dicoding',
            date: new Date(),
            content: 'sebuah comment',
            isDeleted: false,
        };

        // Action
        const modelComment = new ModelComment(payload);

        // Assert
        expect(modelComment).toBeInstanceOf(ModelComment);
        expect(modelComment.id).toStrictEqual(payload.id);
        expect(modelComment.username).toStrictEqual(payload.username);
        expect(modelComment.date).toStrictEqual(payload.date);
        expect(modelComment.content).toStrictEqual(payload.content);
        expect(modelComment.isDeleted).toEqual(payload.isDeleted);
    });

    it('should create Comment object correctly with modified content when isDeleted argument is true', () => {
        // Arrange
        const payload = {
            id: 'comment-123',
            username: 'dicoding',
            date: new Date(),
            content: 'A comment',
            isDeleted: true,
        };

        // Action
        const modelComment = new ModelComment(payload);

        // Assert
        expect(modelComment).toBeInstanceOf(ModelComment);
        expect(modelComment.id).toStrictEqual(payload.id);
        expect(modelComment.username).toStrictEqual(payload.username);
        expect(modelComment.date).toStrictEqual(payload.date);
        expect(modelComment.content).toStrictEqual('**komentar telah dihapus**');
        expect(modelComment.isDeleted).toEqual(payload.isDeleted);
    });
});
