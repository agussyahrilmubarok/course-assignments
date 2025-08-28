export function handleError(error) {
  const response = error.response?.data;

  if (!response) {
    return "Unknown error occurred";
  }

  if (response.errors) {
    return response.errors;
  }

  if (response.message) {
    return response.message;
  }

  return "An unexpected error occurred";
}
