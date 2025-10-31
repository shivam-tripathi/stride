import dotenv from 'dotenv';
import type { ConfigContext, ExpoConfig } from 'expo/config';
import fs from 'node:fs';
import path from 'node:path';

type AppEnv = 'development' | 'production' | 'staging' | 'test';
type PlatformTarget = 'android' | 'ios' | 'web' | 'all';

const APP_NAME = 'react-native-app';
const APP_SLUG = 'react-native-app';

const loadEnvironmentVariables = (appEnv: AppEnv) => {
  const cwd = process.cwd();

  const load = (relativePath: string, override: boolean) => {
    const absolutePath = path.resolve(cwd, relativePath);
    if (!fs.existsSync(absolutePath)) {
      return;
    }

    dotenv.config({ path: absolutePath, override });
  };

  // Load environment variables in order of precedence
  load('.env', false);           // Base configuration
  load(`.env.${appEnv}`, true);  // Environment-specific overrides
  load('.env.local', true);      // Local developer overrides (highest precedence)
};

const toExpoPlatform = (value: string | undefined): PlatformTarget => {
  if (value === 'android' || value === 'ios' || value === 'web') {
    return value;
  }

  return 'all';
};

const resolveAppEnv = (): AppEnv => {
  const explicitEnv = process.env.APP_ENV;
  if (explicitEnv === 'development' || explicitEnv === 'production' || explicitEnv === 'staging' || explicitEnv === 'test') {
    return explicitEnv;
  }

  const fallback = process.env.NODE_ENV === 'production' ? 'production' : 'development';
  return fallback;
};

export default ({ config }: ConfigContext): ExpoConfig => {
  const appEnv = resolveAppEnv();
  const platform = toExpoPlatform(process.env.APP_PLATFORM ?? process.env.EXPO_OS);

  loadEnvironmentVariables(appEnv);

  // Ensure runtime access to the resolved env & platform
  process.env.EXPO_PUBLIC_APP_ENV = appEnv;
  process.env.EXPO_PUBLIC_APP_PLATFORM = platform;

  const baseConfig: ExpoConfig = {
    name: APP_NAME,
    slug: APP_SLUG,
    version: '1.0.0',
    orientation: 'portrait',
    icon: './assets/images/icon.png',
    scheme: 'reactnativeapp',
    userInterfaceStyle: 'automatic',
    newArchEnabled: true,
    ios: {
      supportsTablet: true,
    },
    android: {
      adaptiveIcon: {
        backgroundColor: '#E6F4FE',
        foregroundImage: './assets/images/android-icon-foreground.png',
        backgroundImage: './assets/images/android-icon-background.png',
        monochromeImage: './assets/images/android-icon-monochrome.png',
      },
      edgeToEdgeEnabled: true,
      predictiveBackGestureEnabled: false,
    },
    web: {
      output: 'static',
      favicon: './assets/images/favicon.png',
    },
    plugins: [
      'expo-router',
      [
        'expo-splash-screen',
        {
          image: './assets/images/splash-icon.png',
          imageWidth: 200,
          resizeMode: 'contain',
          backgroundColor: '#ffffff',
          dark: {
            backgroundColor: '#000000',
          },
        },
      ],
    ],
    experiments: {
      typedRoutes: true,
      reactCompiler: true,
    },
  };

  return {
    ...config,
    ...baseConfig,
    extra: {
      ...config.extra,
      ...(baseConfig.extra ?? {}),
      appEnv,
      appPlatform: platform,
    },
  };
};
