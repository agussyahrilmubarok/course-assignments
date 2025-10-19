const Jwt = require('@hapi/jwt');
const config = require('../utils/config');
const AuthenticationError = require('../exceptions/AuthenticationError');

/**
 * This object contains all the methods needed to handle JWT token
 */
const JwtTokenManager = {
  generateAccessToken(payload) {
    // generate the access token using the JWT token generator
    return Jwt.token.generate(payload, config.jwtTokenize.accessTokenKey);
  },

  generateRefreshToken(payload) {
    // generate the refresh token using the JWT token generator
    return Jwt.token.generate(payload, config.jwtTokenize.refreshTokenKey);
  },

  verifyRefreshToken(refreshToken) {
    try {
      // decode the token to get the artifacts (header and payload)
      const artifacts = Jwt.token.decode(refreshToken);
      // verify the signature to check if the token is valid (not tampered)
      Jwt.token.verifySignature(artifacts, config.jwtTokenize.refreshTokenKey);
      // get the decoded payload from the artifacts
      const { payload } = artifacts.decoded;
      // return the payload
      return payload;
    } catch {
      // throw error if token verification failed
      throw new AuthenticationError('refresh token tidak valid');
    }
  },
};

module.exports = JwtTokenManager;
