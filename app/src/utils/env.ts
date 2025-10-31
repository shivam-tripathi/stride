import Constants from 'expo-constants';
import { Platform } from 'react-native';

const extra = Constants.expoConfig?.extra ?? {};

const getExtra = <T>(key: string, fallback: T): T => {
  const value = extra[key];
  return (value as T | undefined) ?? fallback;
};

// Compute appEnv from config (build-time)
const appEnv = getExtra<'development' | 'production' | 'staging' | 'test'>('appEnv', 'development');

// Get platform from config (for reference), but use Platform.OS for runtime detection
const appPlatformConfig = getExtra<'android' | 'ios' | 'web' | 'all'>('appPlatform', 'all');

// Use runtime platform detection from React Native
const appPlatform = Platform.OS as 'android' | 'ios' | 'web';

const getAdjustedApiUrl = (rawUrl: string, platform: string): string => {
  // If not a localhost URL, return as-is
  if (!rawUrl.includes('localhost') && !rawUrl.includes('127.0.0.1')) {
    return rawUrl;
  }

  // Adjust localhost URLs for Android platform
  if (platform === 'android') {
    // Replace localhost/127.0.0.1 with Android emulator host IP
    return rawUrl.replace(/localhost|127\.0\.0\.1/g, '10.0.2.2');
  }

  return rawUrl;
};

const apiUrl = getAdjustedApiUrl(process.env.EXPO_PUBLIC_API_URL ?? 'https://api.example.com', appPlatform);

export const env = {
  appEnv,
  appPlatform,
  apiUrl,
};

export const isDevelopment = appEnv === 'development';
export const isProduction = appEnv === 'production';
export const isStaging = appEnv === 'staging';
export const isTest = appEnv === 'test';

const validateEnvironment = () => {
  const requiredVars = ['EXPO_PUBLIC_API_URL'];

  for (const varName of requiredVars) {
    if (!process.env[varName]) {
      console.warn(`Warning: Required environment variable ${varName} is not set. Using fallback values.`);
    }
  }

  // Validate API URL format
  if (!apiUrl) {
    console.error('Error: API URL is not configured. Please set EXPO_PUBLIC_API_URL in your .env file.');
  } else if (!apiUrl.startsWith('http')) {
    console.warn('Warning: EXPO_PUBLIC_API_URL should be a valid HTTP/HTTPS URL');
  }

  // Validate app environment
  const validEnvs: Array<'development' | 'production' | 'staging' | 'test'> = ['development', 'production', 'staging', 'test'];
  if (!validEnvs.includes(appEnv)) {
    console.warn(`Warning: Invalid APP_ENV '${appEnv}'. Must be one of: ${validEnvs.join(', ')}`);
  }

  // Log detected runtime platform
  console.log(`Runtime platform detected: ${appPlatform}`);

  // Inform about Android localhost adjustment
  const rawApiUrl = process.env.EXPO_PUBLIC_API_URL;
  if (rawApiUrl && (rawApiUrl.includes('localhost') || rawApiUrl.includes('127.0.0.1'))) {
    if (Platform.OS === 'android') {
      console.info(
        `Info: Localhost URL detected for Android. Automatically adjusted to use 10.0.2.2 for Android Emulator.\n` +
        `Original: ${rawApiUrl}\n` +
        `Adjusted: ${apiUrl}\n` +
        `Note: Use your computer's IP address (e.g., 192.168.x.x) for physical devices.`
      );
    }
  }
};

// Validate environment on module load
validateEnvironment();
