/**
 * Redux Store Configuration
 * Combines all reducers and configures the global store
 */

import { configureStore } from '@reduxjs/toolkit';

// Global state reducers
import appReducer from './global/app';
import authReducer from './global/auth';
import themeReducer from './global/theme';

// Feature state reducers
import exploreReducer from './features/explore/exploreSlice';
import homeReducer from './features/home/homeSlice';

export const store = configureStore({
  reducer: {
    // ========================================
    // Global State
    // Accessed from anywhere in the app
    // ========================================
    auth: authReducer,
    theme: themeReducer,
    app: appReducer,

    // ========================================
    // Feature State
    // Primarily used by specific features/screens
    // ========================================
    home: homeReducer,
    explore: exploreReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Ignore these action types if you use non-serializable values
        // ignoredActions: ['your/action/type'],
        // Ignore these paths in the state
        // ignoredPaths: ['items.dates'],
      },
    }),
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

