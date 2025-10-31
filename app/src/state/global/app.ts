/**
 * Global App State
 * Manages app-wide settings and status
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface AppState {
  isOnline: boolean;
  notifications: {
    enabled: boolean;
    count: number;
  };
  lastSync: string | null;
  // Add other app-wide settings here
}

const initialState: AppState = {
  isOnline: true,
  notifications: {
    enabled: true,
    count: 0,
  },
  lastSync: null,
};

const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    setOnlineStatus: (state, action: PayloadAction<boolean>) => {
      state.isOnline = action.payload;
    },
    setNotificationsEnabled: (state, action: PayloadAction<boolean>) => {
      state.notifications.enabled = action.payload;
    },
    setNotificationCount: (state, action: PayloadAction<number>) => {
      state.notifications.count = action.payload;
    },
    incrementNotificationCount: (state) => {
      state.notifications.count += 1;
    },
    clearNotifications: (state) => {
      state.notifications.count = 0;
    },
    setLastSync: (state, action: PayloadAction<string>) => {
      state.lastSync = action.payload;
    },
  },
});

export const {
  setOnlineStatus,
  setNotificationsEnabled,
  setNotificationCount,
  incrementNotificationCount,
  clearNotifications,
  setLastSync,
} = appSlice.actions;
export default appSlice.reducer;

