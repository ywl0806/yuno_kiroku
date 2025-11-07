import { MyAxiosWithAuth } from '../lib/myAxios'

type Range = {
  year: number
  month: number
}

export const getPhotosRange = async (): Promise<Range[]> => {
  const response = await MyAxiosWithAuth.get('/photo/range')

  return response.data
}
