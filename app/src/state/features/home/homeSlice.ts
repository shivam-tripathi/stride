/**
 * Home Feature State
 * Manages state specific to the Home screen
 */

import { apiRequest } from '@/api/client';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';

interface HomeState {
  data: any | null;
  isLoading: boolean;
  error: string | null;
  lastFetched: string | null;
}

const initialState: HomeState = {
  data: null,
  isLoading: false,
  error: null,
  lastFetched: null,
};

// Async thunk for fetching home data
export const fetchHomeData = createAsyncThunk(
  'home/fetchData',
  async (_, { rejectWithValue }) => {
    try {
      // Example endpoint: '/ping' (adjust as needed)
      const response = await apiRequest('/ping');
      return response;
    } catch (error: any) {
      return rejectWithValue(error?.message || 'Failed to fetch home data');
    }
  }
);

const homeSlice = createSlice({
  name: 'home',
  initialState,
  reducers: {
    clearHomeData: (state) => {
      state.data = null;
      state.error = null;
      state.lastFetched = null;
    },
    setHomeError: (state, action) => {
      state.error = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchHomeData.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchHomeData.fulfilled, (state, action) => {
        state.isLoading = false;
        state.data = action.payload;
        state.lastFetched = new Date().toISOString();
        state.error = null;
      })
      .addCase(fetchHomeData.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearHomeData, setHomeError } = homeSlice.actions;
export default homeSlice.reducer;

