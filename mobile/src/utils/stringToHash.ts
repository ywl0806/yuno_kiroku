export const stringToHash = (str: string) => {
  let hash = 0;
  if (str.length === 0) return hash;

  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char; // eslint-disable-line no-bitwise
    // eslint-disable-next-line no-bitwise
    hash &= hash;
  }

  return Math.abs(hash);
};
