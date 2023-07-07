export interface User {
  id?: string;
  name?: string;
  email?: string;
  roles?: Role[];
  token?: string;
}

export interface Role {
  name?: string
}