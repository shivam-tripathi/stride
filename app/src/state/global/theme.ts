/**
 * Global Theme State
 * Manages theme preferences across the app
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

type ColorScheme = 'light' | 'dark' | 'auto';

interface ThemeState {
  colorScheme: ColorScheme;
  // Add other theme-related preferences here
  fontSize: 'small' | 'medium' | 'large';
}

const initialState: ThemeState = {
  colorScheme: 'auto',
  fontSize: 'medium',
};

const themeSlice = createSlice({
  name: 'theme',
  initialState,
  reducers: {
    setColorScheme: (state, action: PayloadAction<ColorScheme>) => {
      state.colorScheme = action.payload;
    },
    setFontSize: (state, action: PayloadAction<ThemeState['fontSize']>) => {
      state.fontSize = action.payload;
    },
    toggleColorScheme: (state) => {
      // Toggle between light and dark (skip auto)
      if (state.colorScheme === 'light') {
        state.colorScheme = 'dark';
      } else if (state.colorScheme === 'dark') {
        state.colorScheme = 'light';
      } else {
        state.colorScheme = 'light';
      }
    },
  },
});

export const { setColorScheme, setFontSize, toggleColorScheme } = themeSlice.actions;
export default themeSlice.reducer;

