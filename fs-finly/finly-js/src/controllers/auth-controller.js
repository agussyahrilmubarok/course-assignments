const bcrypt = require('bcrypt');
const log = require('../libs/logger');
const User = require('../models/user-model');

const { body, validationResult } = require('express-validator');

const validateSignup = [
  body('email', 'Email must not be empty').notEmpty(),
  body('password', 'Password must not be empty').notEmpty(),
  body('password', 'Password must be 6+ characters long').isLength({ min: 6 }),
  body('repeatPassword', 'Repeat Password must not be empty').notEmpty(),
  body('repeatPassword', 'Passwords do not match').custom(
    (value, { req }) => value === req.body.password
  ),
];

const validateLogin = [
  body('email', 'Email must not be empty').notEmpty(),
  body('password', 'Password must not be empty').notEmpty(),
];

const signup = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Signup validation failed', { errors: errors, data: req.body });
    return res.redirect('/signup');
  }

  const { email, password } = req.body;
  const query = { email };

  try {
    const existingUser = await User.findOne(query);
    if (existingUser) {
      req.flash('data', req.body);
      req.flash('info', {
        message: 'Email is already registered. Try to login instead',
        type: 'error',
      });
      log.warn('Signup attempt with existing email', { email });
      return res.redirect('/signup');
    } else {
      const hashedPassword = await bcrypt.hash(password, 10);
      const user = { email, password: hashedPassword };
      const result = await User.create(user);
      req.session.userId = result._id;
      req.flash('info', { message: 'Signup Successful', type: 'success' });
      log.info('New user signed up', {
        email: result.email,
        userId: result._id,
      });
      res.redirect('/dashboard');
    }
  } catch (err) {
    req.flash('info', {
      message: 'An error occurred during signup',
      type: 'error',
    });
    log.error('Error during signup', { error: err.message, stack: err.stack });
    res.redirect('/signup');
  }
};

const login = async (req, res) => {
  const validationErrors = validationResult(req);
  if (!validationErrors.isEmpty()) {
    const errors = validationErrors.array();
    req.flash('errors', errors);
    req.flash('data', req.body);
    log.warn('Login validation failed', { errors: errors, data: req.body });
    return res.redirect('/login');
  }

  const { email, password } = req.body;

  try {
    const user = await User.findOne({ email });
    if (user) {
      const passwordMatch = await bcrypt.compare(password, user.password);
      if (passwordMatch) {
        req.session.userId = user._id;
        req.flash('info', { message: 'Login Successful', type: 'success' });
        log.info('User logged in successfully', {
          email: user.email,
          userId: user._id,
        });
        return res.redirect('/dashboard');
      } else {
        req.flash('info', { message: 'Wrong Password', type: 'error' });
        req.flash('data', req.body);
        log.warn('Login failed - wrong password', { email });
        return res.redirect('/login');
      }
    } else {
      req.flash('info', { message: 'Email is not registered', type: 'error' });
      req.flash('data', req.body);
      log.warn('Login failed - email not found', { email });
      return res.redirect('/login');
    }
  } catch (err) {
    log.error('Error during login', { error: err.message, stack: err.stack });
    return res.redirect('/login');
  }
};

const logout = (req, res) => {
  const userId = req.session.userId;
  if (userId) {
    req.session.userId = null;
    req.flash('info', { message: 'Logout Successful', type: 'success' });
    log.info('User logged out', { userId });
    res.redirect('/');
  } else {
    req.flash('info', { message: 'No active session found', type: 'error' });
    log.warn('Logout attempt with no active session');
    res.redirect('/');
  }
};

module.exports = {
  validateSignup,
  validateLogin,
  signup,
  login,
  logout,
};
