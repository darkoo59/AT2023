import jwt_decode from 'jwt-decode';

export const extractErrorMessages = (res: any) => {
  const err = res.error;
  const messages = [];
  
  if (!err) {
    if (res.message) {
      return [res.message];
    }
    return [];
  }

  if (err.message) messages.push(err.message);

  if (err.title) messages.push(err.title);

  if (err.errors && err.errors.data) messages.push(...err.errors.data);

  return messages;
};

export const deepCopy = (obj: any) => {
  return JSON.parse(JSON.stringify(obj));
};

export const getDecodedAccessToken = (token: string): any => {
  try {
    return jwt_decode(token);
  } catch(Error) {
    return null;
  }
}