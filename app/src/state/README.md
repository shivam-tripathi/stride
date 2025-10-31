# State Management Guide

This project uses **Redux Toolkit** with a **hybrid architecture** that separates global state from feature-specific state.

## Architecture Overview

```
src/state/
├── store.ts              # Redux store configuration
├── hooks.ts              # Typed Redux hooks
├── index.ts              # Central exports
├── global/               # Global state (accessed everywhere)
│   ├── auth.ts          # User authentication
│   ├── theme.ts         # Theme preferences
│   └── app.ts           # App-wide settings
└── features/            # Feature-specific state
    ├── home/
    │   └── homeSlice.ts
    └── explore/
        └── exploreSlice.ts
```

## Philosophy

### Global State (`src/state/global/`)
State that is accessed across many parts of the app:
- ✅ User authentication & profile
- ✅ App theme/preferences
- ✅ Feature flags
- ✅ Notifications
- ✅ Network status

### Feature State (`src/state/features/`)
State that is primarily used by specific features:
- ✅ Feature-specific data fetching
- ✅ Form state for complex forms
- ✅ Filtering, sorting, pagination
- ✅ Feature-specific UI state

### Component Local State
State that should stay in components:
- ✅ UI-only state (dropdowns, modals)
- ✅ Temporary form inputs
- ✅ Animations
- ✅ Focus/hover states

## Usage Examples

### 1. Using Redux in a Component

```typescript
import { useAppDispatch, useAppSelector } from '@/state/hooks';
import { fetchHomeData } from '@/state/features/home/homeSlice';
import { logout } from '@/state/global/auth';

function MyComponent() {
  const dispatch = useAppDispatch();

  // Access feature state
  const { data, isLoading } = useAppSelector((state) => state.home);

  // Access global state
  const { user } = useAppSelector((state) => state.auth);

  // Dispatch actions
  const handleFetch = () => {
    dispatch(fetchHomeData());
  };

  const handleLogout = () => {
    dispatch(logout());
  };

  return (
    <View>
      <Text>Welcome, {user?.name}</Text>
      <Button onPress={handleFetch} loading={isLoading}>
        Fetch Data
      </Button>
    </View>
  );
}
```

### 2. Creating a New Feature Slice

```typescript
// src/state/features/profile/profileSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiRequest } from '@/api/client';

interface ProfileState {
  profile: UserProfile | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ProfileState = {
  profile: null,
  isLoading: false,
  error: null,
};

// Async thunk for API calls
export const fetchProfile = createAsyncThunk(
  'profile/fetch',
  async (userId: string, { rejectWithValue }) => {
    try {
      return await apiRequest(`/users/${userId}`);
    } catch (error: any) {
      return rejectWithValue(error?.message || 'Failed to fetch profile');
    }
  }
);

const profileSlice = createSlice({
  name: 'profile',
  initialState,
  reducers: {
    clearProfile: (state) => {
      state.profile = null;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchProfile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.profile = action.payload;
      })
      .addCase(fetchProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearProfile } = profileSlice.actions;
export default profileSlice.reducer;
```

**Then add it to the store:**

```typescript
// src/state/store.ts
import profileReducer from './features/profile/profileSlice';

export const store = configureStore({
  reducer: {
    // ... other reducers
    profile: profileReducer,
  },
});
```

### 3. Accessing State with Selectors

For better performance and reusability, create selector functions:

```typescript
// src/state/features/home/selectors.ts
import { RootState } from '@/state/store';

export const selectHomeData = (state: RootState) => state.home.data;
export const selectHomeLoading = (state: RootState) => state.home.isLoading;
export const selectHomeError = (state: RootState) => state.home.error;

// Derived selector
export const selectHasHomeData = (state: RootState) =>
  state.home.data !== null && !state.home.error;

// Usage in component
import { selectHomeData, selectHasHomeData } from '@/state/features/home/selectors';

const data = useAppSelector(selectHomeData);
const hasData = useAppSelector(selectHasHomeData);
```

### 4. Complex Async Operations

```typescript
// Thunk with multiple API calls
export const updateProfileWithAvatar = createAsyncThunk(
  'profile/updateWithAvatar',
  async ({ profileData, avatarFile }: UpdateProfileArgs, { dispatch }) => {
    // Upload avatar first
    const avatarUrl = await apiRequest('/upload', {
      method: 'POST',
      body: avatarFile,
    });

    // Update profile with avatar URL
    const updatedProfile = await apiRequest('/profile', {
      method: 'PUT',
      body: JSON.stringify({ ...profileData, avatar: avatarUrl }),
    });

    // Update auth user if needed
    dispatch(setUser(updatedProfile));

    return updatedProfile;
  }
);
```

## Best Practices

### ✅ Do's

1. **Use typed hooks** (`useAppDispatch`, `useAppSelector`) everywhere
2. **Keep feature state isolated** - each feature has its own slice
3. **Use `createAsyncThunk`** for all async operations
4. **Handle loading and error states** in every slice
5. **Export actions** from slice files for easy imports
6. **Use selectors** for computed/derived state
7. **Keep slices focused** - one slice per feature or domain

### ❌ Don'ts

1. **Don't put everything in Redux** - use local state when appropriate
2. **Don't mutate state directly** - Redux Toolkit handles this with Immer
3. **Don't store derived data** - compute it with selectors
4. **Don't duplicate state** between slices
5. **Don't store non-serializable values** (functions, promises, class instances)
6. **Don't make massive slices** - split large features into sub-features

## Redux DevTools

Redux DevTools is automatically enabled in development mode. You can:
- Time-travel through state changes
- Inspect actions and payloads
- See state diffs
- Export/import state snapshots

## Migration Guide

If you need to migrate existing component state to Redux:

1. **Identify the scope**: Is it global or feature-specific?
2. **Create/update the slice**: Add necessary state and reducers
3. **Replace useState with useAppSelector**: Access state from Redux
4. **Replace setState with dispatch**: Dispatch actions instead
5. **Test thoroughly**: Ensure all functionality works

## Performance Tips

1. **Use `reselect`** for expensive computations:
```typescript
import { createSelector } from '@reduxjs/toolkit';

const selectExpensiveData = createSelector(
  [(state) => state.data.items],
  (items) => items.filter(/* expensive operation */)
);
```

2. **Split selectors** - avoid selecting entire slices:
```typescript
// ❌ Bad - causes re-render on any home state change
const home = useAppSelector(state => state.home);

// ✅ Good - only re-renders when data changes
const data = useAppSelector(state => state.home.data);
```

3. **Normalize nested data** with `@reduxjs/toolkit/query` or normalized state shape

## Resources

- [Redux Toolkit Docs](https://redux-toolkit.js.org/)
- [Redux Style Guide](https://redux.js.org/style-guide/)
- [Redux DevTools](https://github.com/reduxjs/redux-devtools)

