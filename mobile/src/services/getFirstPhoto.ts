import {MyAxios} from '../lib/myAxios';
import {Photo} from '../types/photo';

export const getFirstPhoto = async (): Promise<Photo> => {
  const url = `http://localhost:1323/api/photo/first`;
  const response = await MyAxios.get(url);
  return response.data;
};
