const CurrencyUSDollar = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
});

const CurrencyCNY = new Intl.NumberFormat('zh-CN', {
  style: 'currency',
  currency: 'CNY',
});

const CurrencyJPY = new Intl.NumberFormat('ja-JP', {
  style: 'currency',
  currency: 'JPY',
});

const CurrencyIDR = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
});

module.exports = {
  CurrencyUSDollar,
  CurrencyCNY,
  CurrencyJPY,
  CurrencyIDR,
};
