declare global {
  namespace NodeJS {
    interface ProcessEnv {
      EXPO_PUBLIC_API_URL?: string;
      EXPO_PUBLIC_APP_ENV?: 'development' | 'production' | 'staging' | 'test';
      EXPO_PUBLIC_APP_PLATFORM?: 'android' | 'ios' | 'web' | 'all';
      APP_ENV?: 'development' | 'production' | 'staging' | 'test';
      APP_PLATFORM?: string;
      EXPO_OS?: string;
      NODE_ENV?: string;
    }
  }
}

export { };

