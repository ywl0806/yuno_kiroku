import {NativeModules} from 'react-native';

const {LivePhotoModule} = NativeModules;

export const getLivePhotoMov = (id: string): Promise<string> => {
  return new Promise((resolve, reject) => {
    LivePhotoModule.getLivePhotoMov(id)
      .then((r: string) => {
        resolve(r);
      })
      .catch((e: any) => {
        reject(e);
      });
  });
};
