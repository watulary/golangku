export function validatePassword(password) {
  return {
    minLength: password.length >= 8,
    maxLength: password.length <= 64,
    uppercase: /[A-Z]/.test(password),
    lowercase: /[a-z]/.test(password),
    number: /[0-9]/.test(password),
    special: /[^A-Za-z0-9]/.test(password),
  };
}

export function isPasswordValid(password) {
  const result = validatePassword(password);

  return Object.values(result).every(Boolean);
}
