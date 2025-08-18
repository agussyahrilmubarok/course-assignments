export class User {
  constructor({ id, fullName, email, roles }) {
    this.id = id;
    this.fullName = fullName;
    this.email = email;
    this.roles = roles || [];
  }

  isAdmin() {
    return this.roles.includes("ROLE_ADMIN");
  }

  isUser() {
    return this.roles.includes("ROLE_USER");
  }
}
