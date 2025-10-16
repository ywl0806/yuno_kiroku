import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';

// This function uploads a Live Photo by sending both the photo and live video files to the server as form data.

export const uploadLivePhoto = async (photo: PhotoIdentifier, live: string) => {
  // Create a new instance of FormData to hold the files
  const formData = new FormData();

  if (!photo.node.image.filename) {
    return;
  }

  // Split the filename and extension of the photo file
  const [filename, ext] = photo.node.image.filename.split('.');

  // Append the photo file to the form data
  formData.append('photo', {
    uri: photo.node.image.uri,
    type: 'image/*',
    name: filename + '.' + ext,
  });

  // Append the live video file to the form data
  formData.append('live', {
    uri: live,
    type: 'video/*',
    name: filename + '.mov',
  });

  // Send a POST request to the server with the form data
  const response = await fetch('http://localhost:1323/api/photo/upload-live', {
    method: 'POST',
    body: formData,
  });

  // Parse the response as JSON and return it
  return response.json();
};
