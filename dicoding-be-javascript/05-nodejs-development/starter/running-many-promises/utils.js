// DON'T CHANGE THE CODE BELOW

function getUserInfo(userId) {
  return new Promise((resolve, reject) => {
    if (userId === -1) {
      reject(new Error('Invalid user id'));
    }

    setTimeout(() => {
      resolve({ id: userId, name: 'Dicoding' });
    }, 1000);
  });
}

function getUserOrders(userId) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([{ order: 1, userId }, { order: 2, userId }, { order: 3, userId }]);
    }, 500);
  });
}

function getUserCartItems(userId) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([{ item: 1, userId }, { item: 2, userId }]);
    }, 250);
  });
}

module.exports = { getUserInfo, getUserOrders, getUserCartItems };
