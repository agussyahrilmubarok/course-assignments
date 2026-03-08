function f(n) {
  if (n === 15) {
    throw new Error('Error!');
  }

  return f(n + 1);
}

f(1);
