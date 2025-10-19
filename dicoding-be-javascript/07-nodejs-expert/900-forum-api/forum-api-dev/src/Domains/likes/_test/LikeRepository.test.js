const LikeRepository = require('../LikeRepository');

describe('a LikeRepository interface', () => {
    it('should throw error when invoke abstract behavior', async () => {
        // Arrange
        const likeRepository = new LikeRepository();

        // Action & Assert
        await expect(likeRepository.addLikeToComment({})).rejects.toThrowError('LIKE_REPOSITORY.METHOD_NOT_IMPLEMENTED');
        await expect(likeRepository.getLikeCountComment({})).rejects.toThrowError('LIKE_REPOSITORY.METHOD_NOT_IMPLEMENTED');
        await expect(likeRepository.verifyLikeIsExists({})).rejects.toThrowError('LIKE_REPOSITORY.METHOD_NOT_IMPLEMENTED');
        await expect(likeRepository.unLikeComment({})).rejects.toThrowError('LIKE_REPOSITORY.METHOD_NOT_IMPLEMENTED');
    });
});
