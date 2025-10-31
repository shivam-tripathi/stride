/**
 * React Native Paper theme configuration
 * Integrates with the app's existing color scheme
 */

import type { MD3Theme } from 'react-native-paper';
import { MD3DarkTheme, MD3LightTheme } from 'react-native-paper';
import { Colors } from './colors';

export const paperLightTheme: MD3Theme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    primary: Colors.light.tint,
    onPrimary: '#fff',
    primaryContainer: '#d1e4ff',
    onPrimaryContainer: '#001d36',
    secondary: '#535f70',
    onSecondary: '#fff',
    secondaryContainer: '#d7e3f7',
    onSecondaryContainer: '#101c2b',
    background: Colors.light.background,
    onBackground: Colors.light.text,
    surface: Colors.light.background,
    onSurface: Colors.light.text,
    surfaceVariant: '#dfe2eb',
    onSurfaceVariant: '#43474e',
    outline: '#73777f',
    outlineVariant: '#c3c7cf',
  },
};

export const paperDarkTheme: MD3Theme = {
  ...MD3DarkTheme,
  colors: {
    ...MD3DarkTheme.colors,
    primary: Colors.dark.tint,
    onPrimary: '#003258',
    primaryContainer: '#00497d',
    onPrimaryContainer: '#d1e4ff',
    secondary: '#bbc7db',
    onSecondary: '#253140',
    secondaryContainer: '#3b4858',
    onSecondaryContainer: '#d7e3f7',
    background: Colors.dark.background,
    onBackground: Colors.dark.text,
    surface: Colors.dark.background,
    onSurface: Colors.dark.text,
    surfaceVariant: '#43474e',
    onSurfaceVariant: '#c3c7cf',
    outline: '#8d9199',
    outlineVariant: '#43474e',
  },
};

