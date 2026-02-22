import { PermissionsAndroid, Platform } from 'react-native';

export const requestAndroidPermissions = async (): Promise<{
  phone: boolean;
  notifications: boolean;
}> => {
  const results = {
    phone: false,
    notifications: false,
  };

  if (Platform.OS !== 'android') {
    return results;
  }

  try {
    if (Platform.OS === 'android') {
      const androidVersion = Platform.Version;

      if (androidVersion >= 33) {
        const notificationGranted = await PermissionsAndroid.request(
          PermissionsAndroid.PERMISSIONS.POST_NOTIFICATIONS,
          {
            title: 'Notification Permission',
            message: 'Eco needs access to notifications to sync them to your PC.',
            buttonNeutral: 'Ask Me Later',
            buttonNegative: 'Cancel',
            buttonPositive: 'OK',
          }
        );
        results.notifications = notificationGranted === PermissionsAndroid.RESULTS.GRANTED;
      }

      const phoneGranted = await PermissionsAndroid.request(
        PermissionsAndroid.PERMISSIONS.READ_PHONE_STATE,
        {
          title: 'Phone State Permission',
          message: 'Eco needs access to phone state to detect incoming calls.',
          buttonNeutral: 'Ask Me Later',
          buttonNegative: 'Cancel',
          buttonPositive: 'OK',
        }
      );
      results.phone = phoneGranted === PermissionsAndroid.RESULTS.GRANTED;
    }
  } catch (err) {
    console.warn('Permission request error:', err);
  }

  console.log('Permission results:', results);
  return results;
};

export const checkAndroidPermissions = async (): Promise<{
  phone: boolean;
  notifications: boolean;
}> => {
  const results = {
    phone: false,
    notifications: false,
  };

  if (Platform.OS !== 'android') {
    return results;
  }

  try {
    if (Platform.OS === 'android') {
      const androidVersion = Platform.Version;

      if (androidVersion >= 33) {
        results.notifications = await PermissionsAndroid.check(
          PermissionsAndroid.PERMISSIONS.POST_NOTIFICATIONS
        );
      }

      results.phone = await PermissionsAndroid.check(
        PermissionsAndroid.PERMISSIONS.READ_PHONE_STATE
      );
    }
  } catch (err) {
    console.warn('Permission check error:', err);
  }

  return results;
};
