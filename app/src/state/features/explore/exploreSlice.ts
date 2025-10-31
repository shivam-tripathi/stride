/**
 * Explore Feature State
 * Manages state specific to the Explore screen
 */

import { apiRequest } from '@/api/client';
import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

interface ExploreItem {
  id: string;
  title: string;
  description: string;
  imageUrl?: string;
}

interface ExploreState {
  items: ExploreItem[];
  isLoading: boolean;
  error: string | null;
  searchQuery: string;
  selectedCategory: string | null;
  page: number;
  hasMore: boolean;
}

const initialState: ExploreState = {
  items: [],
  isLoading: false,
  error: null,
  searchQuery: '',
  selectedCategory: null,
  page: 1,
  hasMore: true,
};

// Async thunk for fetching explore items
export const fetchExploreItems = createAsyncThunk(
  'explore/fetchItems',
  async ({ page, category }: { page: number; category?: string | null }, { rejectWithValue }) => {
    try {
      const endpoint = category
        ? `/explore?page=${page}&category=${category}`
        : `/explore?page=${page}`;
      const response = await apiRequest<ExploreItem[]>(endpoint);
      return response;
    } catch (error: any) {
      return rejectWithValue(error?.message || 'Failed to fetch explore items');
    }
  }
);

const exploreSlice = createSlice({
  name: 'explore',
  initialState,
  reducers: {
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
    setSelectedCategory: (state, action: PayloadAction<string | null>) => {
      state.selectedCategory = action.payload;
      state.page = 1;
      state.items = [];
      state.hasMore = true;
    },
    clearExploreItems: (state) => {
      state.items = [];
      state.error = null;
      state.page = 1;
      state.hasMore = true;
    },
    incrementPage: (state) => {
      state.page += 1;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchExploreItems.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchExploreItems.fulfilled, (state, action) => {
        state.isLoading = false;
        // Append new items for pagination
        state.items = [...state.items, ...action.payload];
        state.hasMore = action.payload.length > 0;
        state.error = null;
      })
      .addCase(fetchExploreItems.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { setSearchQuery, setSelectedCategory, clearExploreItems, incrementPage } = exploreSlice.actions;
export default exploreSlice.reducer;

