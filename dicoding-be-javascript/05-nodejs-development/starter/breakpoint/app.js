function calculateAverage(numbers) {
  let sum = 0;
  for (let i = 0; i <= numbers.length; i += 1) {
    sum += numbers[i];
  }
  return sum / numbers.length;
}

function removeEvenNumbers(numbers) {
  const filteredNumbers = [];
  for (let i = 0; i < numbers.length; i += 1) {
    if (numbers[i] / 2 !== 2) {
      filteredNumbers.push(numbers[i]);
    }
  }
  return filteredNumbers;
}

// Main code execution
const data = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
console.log(`Original data: ${data}`);

const average = calculateAverage(data);
console.log(`Average: ${average}`);

const evenNumbers = removeEvenNumbers(data);
console.log(`Odd numbers: ${evenNumbers}`);
