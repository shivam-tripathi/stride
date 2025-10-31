# React Native App

A scalable React Native application built with Expo, featuring a clean architecture that separates routing from application logic.


## Quick Start

```bash
# Install dependencies
npm install

# Start the development server (default: dev mode, platform auto-detect)
npx expo start
```

## Environment Variables & Build Modes

This project supports environment variables for different build types (development, production, staging, test) using `.env` files and a dynamic Expo config.

### How it works

- Environment variables are loaded from `.env`, `.env.[env]`, and `.env.local` (in that order of precedence).
- Only variables prefixed with `EXPO_PUBLIC_` are available at runtime in the app.
- See `.env.example` and other `.env.*.example` files for templates.
- The system validates required variables and provides helpful warnings with detailed error messages.
- TypeScript type definitions ensure type safety for all environment variables (see `src/types/env.d.ts`).

### Example: Setting API URL

```
EXPO_PUBLIC_API_URL=https://api.example.com
```

### Running for Different Environments

#### Android
```bash
# Development
npm run android:dev
# Production
npm run android:prod
# Staging
npm run android:staging
```

#### iOS
```bash
# Development
npm run ios:dev
# Production
npm run ios:prod
# Staging
npm run ios:staging
```

#### Web
```bash
# Development
npm run web:dev
# Production
npm run web:prod
# Staging
npm run web:staging
```

#### Test Environment
```bash
npm run test
```

You can also set custom environment:
```bash
cross-env APP_ENV=staging expo start --web
```

#### Adding/Editing Environment Variables
- Copy `.env.example` to `.env` and adjust values for your project.
- For environment-specific overrides, copy `.env.development.example`, `.env.production.example`, etc.
- For local development overrides (never committed), create `.env.local` with your personal settings.
- **Android Development Note**: The system automatically adjusts `localhost` URLs to `10.0.2.2` for Android emulators. You'll see an info message in the console when this happens. For physical Android devices, use your computer's IP address (e.g., `192.168.x.x`).
- Never commit real `.env` files—only the `.example` templates are tracked.

#### Available Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `EXPO_PUBLIC_API_URL` | Backend API URL | Yes | `https://api.example.com` |
| `EXPO_PUBLIC_APP_ENV` | App environment | No | Auto-set by scripts |
| `EXPO_PUBLIC_APP_PLATFORM` | Target platform | No | Auto-detected |

The app will validate these on startup and show clear error messages if configuration is incorrect.

#### Accessing Variables in Code
Use the centralized helper in `src/utils/env.ts`:

```typescript
import { env, isDevelopment, isProduction, isStaging, isTest } from '@/utils/env';

console.log(env.apiUrl); // Reads EXPO_PUBLIC_API_URL (with platform adjustments)
console.log(env.appEnv); // Current environment: 'development' | 'production' | 'staging' | 'test'
console.log(env.appPlatform); // Current platform: 'android' | 'ios' | 'web' | 'all'

// Environment checks
if (isDevelopment) {
  // Development-specific code
}

if (isProduction) {
  // Production-specific code
}
```

**Important**: Always use `env.apiUrl` from `@/utils/env` for API requests. This ensures proper platform-specific URL handling (e.g., Android emulator localhost adjustment) and provides a single source of truth for configuration.

---

## UI Components & Theming

This project uses **React Native Paper** for Material Design components with a custom theme that adapts to light/dark mode.

### React Native Paper

React Native Paper provides a comprehensive set of Material Design components that are:
- Accessible and production-ready
- Fully themeable with automatic dark mode support
- Consistent across platforms (iOS, Android, Web)

**Available Components:**
- `Button`, `Card`, `Text`, `Title`, `Paragraph`
- `TextInput`, `FAB`, `Chip`, `Badge`, `Avatar`
- `Dialog`, `Modal`, `Snackbar`, `Banner`
- `List`, `DataTable`, `Divider`, `Surface`
- And many more! See [React Native Paper docs](https://callstack.github.io/react-native-paper/)

**Usage Example:**
```typescript
import { Button, Card, Text } from 'react-native-paper';

<Card>
  <Card.Content>
    <Text variant="titleLarge">Hello World</Text>
    <Button mode="contained" onPress={handlePress}>
      Press me
    </Button>
  </Card.Content>
</Card>
```

### Theme Customization

The theme is configured in `src/theme/paper-theme.ts` and automatically syncs with your app's color scheme:

```typescript
import { paperLightTheme, paperDarkTheme } from '@/theme';
import { useTheme } from 'react-native-paper';

// In your component
const theme = useTheme();
const primaryColor = theme.colors.primary;
```

To customize the theme, edit `src/theme/paper-theme.ts` and `src/theme/colors.ts`.

---

## State Management

This project uses **Redux Toolkit** with a **hybrid architecture** that separates global state from feature-specific state.

### Architecture

```
src/state/
├── store.ts              # Redux store configuration
├── hooks.ts              # Typed Redux hooks (useAppDispatch, useAppSelector)
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

### When to Use Redux vs Local State

**Global State (`src/state/global/`)** - Accessed across many parts of the app:
- ✅ User authentication & profile
- ✅ App theme/preferences
- ✅ Feature flags
- ✅ Notifications
- ✅ Network status

**Feature State (`src/state/features/`)** - Primarily used by specific features:
- ✅ Feature-specific data fetching
- ✅ Form state for complex forms
- ✅ Filtering, sorting, pagination
- ✅ Feature-specific UI state

**Component Local State** - Stay in components:
- ✅ UI-only state (dropdowns, modals)
- ✅ Temporary form inputs
- ✅ Animations

### Usage Example

```typescript
import { useAppDispatch, useAppSelector } from '@/state/hooks';
import { fetchHomeData } from '@/state/features/home/homeSlice';
import { logout } from '@/state/global/auth';

function HomeScreen() {
  const dispatch = useAppDispatch();

  // Access feature-specific state
  const { data, isLoading, error } = useAppSelector((state) => state.home);

  // Access global state
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);

  const handleFetch = () => {
    dispatch(fetchHomeData());
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

See the [State Management Guide](src/state/README.md) for detailed documentation, best practices, and examples.

---

## Architecture

This project uses a **scalable architecture** that separates routing concerns from application logic, making the codebase easier to navigate, test, and maintain.

### Directory Structure

```
├── app/                    # Routing layer (Expo Router)
│   ├── _layout.tsx
│   ├── (tabs)/
│   │   ├── _layout.tsx
│   │   ├── index.tsx       # → HomeScreen
│   │   └── explore.tsx     # → ExploreScreen
│   └── modal.tsx           # → ModalScreen
│
├── src/                    # Application code
│   ├── api/                # API clients and services
│   ├── assets/             # App-specific assets (fonts, images)
│   ├── components/         # Reusable React components
│   │   ├── common/         # Simple UI components
│   │   └── ui/             # Complex UI elements
│   ├── constants/          # App-wide constants
│   ├── hooks/              # Custom React hooks
│   ├── i18n/               # Internationalization
│   ├── screens/            # Screen components
│   ├── state/              # Redux state management
│   │   ├── global/         # Global state (auth, theme, app)
│   │   ├── features/       # Feature-specific state
│   │   ├── store.ts        # Redux store configuration
│   │   └── hooks.ts        # Typed Redux hooks
│   ├── theme/              # Design system (colors, fonts, metrics)
│   ├── types/              # TypeScript types
│   └── utils/              # Utility functions
│
└── assets/                 # Native app assets (icons, splash)
```

### Core Principles

**Separation of Concerns**
- `app/` - Routing only
- `src/` - All application logic
- Clear boundaries, predictable structure

**Scalability**
- Easy to find where code belongs
- Each directory has a single purpose
- Adding features doesn't bloat routing

**Maintainability**
- Consistent patterns
- Easy onboarding
- Reduced coupling

**Path Aliases**
- Use `@/*` to import from `src/`
- Clean imports: `@/components/Button` vs `../../../components/Button`

## Development Guide

### Adding a New Screen

1. Create screen component in `src/screens/[feature]/`
2. Create route file in `app/` that imports the screen
3. Update navigation if needed

**Example:**

```typescript
// src/screens/profile/ProfileScreen.tsx
import { StyleSheet } from 'react-native';
import { Surface, Text } from 'react-native-paper';
import { Spacing } from '@/theme';

export default function ProfileScreen() {
  return (
    <Surface style={styles.container}>
      <Text variant="headlineMedium">Profile</Text>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: Spacing.md
  },
});
```

```typescript
// app/profile.tsx
import ProfileScreen from '@/screens/profile/ProfileScreen';
export default ProfileScreen;
```

### Creating Reusable Components

Place in `src/components/common/` or `src/components/ui/`. **Always use React Native Paper components** for consistent Material Design:

```typescript
// src/components/common/UserCard.tsx
import { Card, Avatar, Button } from 'react-native-paper';
import { StyleSheet } from 'react-native';
import { Spacing } from '@/theme';

type UserCardProps = {
  name: string;
  email: string;
  avatarUrl?: string;
  onPress: () => void;
};

export function UserCard({ name, email, avatarUrl, onPress }: UserCardProps) {
  return (
    <Card style={styles.card} mode="elevated">
      <Card.Title
        title={name}
        subtitle={email}
        left={(props) => <Avatar.Image {...props} source={{ uri: avatarUrl }} />}
      />
      <Card.Actions>
        <Button onPress={onPress}>View Profile</Button>
      </Card.Actions>
    </Card>
  );
}

const styles = StyleSheet.create({
  card: {
    marginVertical: Spacing.sm,
  },
});
```

**Accessing Theme Values:**

```typescript
import { Spacing, BorderRadius } from '@/theme';
import { useTheme } from 'react-native-paper';
import { StyleSheet } from 'react-native';

function MyComponent() {
  const theme = useTheme();

  const styles = StyleSheet.create({
    container: {
      padding: Spacing.md,
      borderRadius: BorderRadius.md,
      backgroundColor: theme.colors.surface,
    },
  });

  return <Surface style={styles.container}>...</Surface>;
}
```

### Creating Custom Hooks

Place in `src/hooks/`

```typescript
// src/hooks/use-toggle.ts
import { useState } from 'react';

export function useToggle(initialValue = false) {
  const [value, setValue] = useState(initialValue);
  const toggle = () => setValue(v => !v);
  return [value, toggle] as const;
}
```

### Using the Theme System

The theme system consists of:
- **Paper Theme** (`src/theme/paper-theme.ts`) - Material Design colors for Paper components
- **Custom Constants** (`src/theme/`) - Spacing, fonts, metrics for layout

**Customize Paper Theme Colors:**

```typescript
// src/theme/colors.ts
export const Colors = {
  light: {
    text: '#11181C',
    background: '#fff',
    tint: '#0a7ea4',  // This becomes Paper's primary color
    // ... other colors
  },
  dark: {
    text: '#ECEDEE',
    background: '#151718',
    tint: '#fff',
    // ... other colors
  },
};
```

**Use in Components:**

```typescript
import { Spacing, BorderRadius } from '@/theme';
import { useTheme } from 'react-native-paper';
import { Surface, Text } from 'react-native-paper';
import { StyleSheet } from 'react-native';

function MyScreen() {
  const theme = useTheme();

  return (
    <Surface style={styles.container}>
      <Text variant="headlineMedium" style={{ color: theme.colors.primary }}>
        Hello
      </Text>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: Spacing.lg,
    borderRadius: BorderRadius.md,
  },
});
```

### Adding API Services

Place in `src/api/`. Always use the centralized `env.apiUrl` for API configuration:

```typescript
// src/api/client.ts
import { env } from '@/utils/env';

export async function apiRequest<T = any>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${env.apiUrl}${endpoint}`, options);
  if (!response.ok) throw new Error('API Error');
  return response.json();
}
```

```typescript
// src/api/services/user.ts
import { apiRequest } from '@/api/client';

export const userService = {
  getProfile: () => apiRequest('/user/profile'),
  updateProfile: (data: any) =>
    apiRequest('/user/profile', { method: 'PUT', body: JSON.stringify(data) }),
};
```

### Adding Utility Functions

Place in `src/utils/`

```typescript
// src/utils/format.ts
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(amount);
}

export function truncate(str: string, length: number): string {
  return str.length > length ? str.slice(0, length) + '...' : str;
}
```

### Adding TypeScript Types

Place in `src/types/`

```typescript
// src/types/models.ts
export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
}

export interface Post {
  id: string;
  title: string;
  content: string;
  authorId: string;
  createdAt: Date;
}
```

### State Management

Redux is already configured! See the [State Management Guide](src/state/README.md) for complete documentation.

**Quick example - Adding a new feature slice:**

```typescript
// src/state/features/profile/profileSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiRequest } from '@/api/client';

export const fetchProfile = createAsyncThunk(
  'profile/fetch',
  async (userId: string) => {
    return await apiRequest(`/users/${userId}`);
  }
);

const profileSlice = createSlice({
  name: 'profile',
  initialState: { data: null, isLoading: false, error: null },
  reducers: {
    clearProfile: (state) => {
      state.data = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchProfile.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.data = action.payload;
      })
      .addCase(fetchProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message;
      });
  },
});

export const { clearProfile } = profileSlice.actions;
export default profileSlice.reducer;
```

Then add to `src/state/store.ts`:
```typescript
import profileReducer from './features/profile/profileSlice';

export const store = configureStore({
  reducer: {
    // ...
    profile: profileReducer,
  },
});
```

## Best Practices

- ✅ Keep route files minimal (import/export only)
- ✅ **Always use React Native Paper components** (Button, Card, Text, Surface, etc.)
- ✅ Use Paper's `useTheme()` for accessing theme colors
- ✅ Use custom theme constants (`Spacing`, `BorderRadius`) for layout
- ✅ No hardcoded spacing values - use `Spacing.md`, `Spacing.lg`, etc.
- ✅ Organize by feature (related files together)
- ✅ Extract reusable logic into components/hooks
- ✅ Use TypeScript types consistently
- ✅ Use path aliases (`@/*`) for cleaner imports

## Project Structure Benefits

| Aspect | Benefit |
|--------|---------|
| **Finding code** | Each directory has a clear purpose |
| **Adding features** | Obvious where new code belongs |
| **Testing** | Logic decoupled from routing |
| **Onboarding** | Predictable, consistent patterns |
| **Refactoring** | Changes isolated to specific areas |

---

Built with [Expo](https://expo.dev) and [React Native](https://reactnative.dev)
