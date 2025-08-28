import test from 'node:test';
import assert from 'node:assert';
import sum from './index.js';

// Case 1: kedua parameter valid number positif
test('sum of two positive numbers', () => {
  assert.strictEqual(sum(2, 3), 5);
  assert.strictEqual(sum(10, 20), 30);
});

// Case 2: salah satu parameter bukan number
test('return 0 if one parameter is not a number', () => {
  assert.strictEqual(sum('2', 3), 0);
  assert.strictEqual(sum(2, '3'), 0);
});

// Case 3: kedua parameter bukan number
test('return 0 if both parameters are not numbers', () => {
  assert.strictEqual(sum('a', 'b'), 0);
});

// Case 4: salah satu parameter negatif
test('return 0 if one parameter is negative', () => {
  assert.strictEqual(sum(-1, 5), 0);
  assert.strictEqual(sum(5, -1), 0);
});

// Case 5: kedua parameter negatif
test('return 0 if both parameters are negative', () => {
  assert.strictEqual(sum(-2, -3), 0);
});

// Case 6: kedua parameter nol
test('sum of zeros should be zero', () => {
  assert.strictEqual(sum(0, 0), 0);
});
