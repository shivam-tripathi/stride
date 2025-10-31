# Theme System Guide

This project uses **React Native Paper** for UI components with a custom theme configuration that provides a single source of truth for colors, spacing, and other design tokens.

## Architecture

```
src/theme/
├── colors.ts           # Base color palette (feeds into Paper theme)
├── paper-theme.ts      # Material Design 3 theme for Paper components
├── fonts.ts            # Platform-specific font configurations
├── metrics.ts          # Spacing, BorderRadius, IconSize constants
└── index.ts            # Central export
```

## How It Works

### 1. Base Colors → Paper Theme

Define colors once in `colors.ts`, and they automatically flow into Paper's theme:

```typescript
// src/theme/colors.ts
export const Colors = {
  light: {
    tint: '#0a7ea4',      // → Paper's primary color
    background: '#fff',    // → Paper's background
    text: '#11181C',       // → Paper's onBackground
  },
  dark: { /* ... */ }
};
```

```typescript
// src/theme/paper-theme.ts
export const paperLightTheme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    primary: Colors.light.tint,           // ← Uses your color
    background: Colors.light.background,  // ← Uses your color
    onBackground: Colors.light.text,      // ← Uses your color
  },
};
```

### 2. Using Paper Components

**Always use Paper components** for UI elements:

```typescript
import { Surface, Text, Button, Card } from 'react-native-paper';
import { Spacing } from '@/theme';

function MyScreen() {
  return (
    <Surface style={{ padding: Spacing.md }}>
      <Card>
        <Card.Content>
          <Text variant="headlineMedium">Title</Text>
          <Text variant="bodyMedium">Description</Text>
          <Button mode="contained">Action</Button>
        </Card.Content>
      </Card>
    </Surface>
  );
}
```

### 3. Accessing Theme in Components

Use Paper's `useTheme()` hook for dynamic theme values:

```typescript
import { useTheme } from 'react-native-paper';
import { Surface, Text } from 'react-native-paper';
import { Spacing, BorderRadius } from '@/theme';
import { StyleSheet } from 'react-native';

function MyComponent() {
  const theme = useTheme();

  return (
    <Surface style={styles.container}>
      <Text style={{ color: theme.colors.primary }}>
        Dynamic color!
      </Text>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: Spacing.md,
    borderRadius: BorderRadius.lg,
  },
});
```

## Available Design Tokens

### Colors

Colors are managed by Paper. Access them via `useTheme()`:

```typescript
const theme = useTheme();

// Available colors (Material Design 3)
theme.colors.primary
theme.colors.secondary
theme.colors.background
theme.colors.surface
theme.colors.error
theme.colors.onPrimary
theme.colors.onBackground
// ... and many more
```

### Spacing

Consistent spacing across the app:

```typescript
import { Spacing } from '@/theme';

Spacing.xs   // 4
Spacing.sm   // 8
Spacing.md   // 16
Spacing.lg   // 24
Spacing.xl   // 32
Spacing.xxl  // 48
```

**Always use these instead of hardcoded values!**

```typescript
// ❌ Bad
style={{ padding: 16 }}

// ✅ Good
style={{ padding: Spacing.md }}
```

### Border Radius

```typescript
import { BorderRadius } from '@/theme';

BorderRadius.sm    // 4
BorderRadius.md    // 8
BorderRadius.lg    // 16
BorderRadius.xl    // 24
BorderRadius.full  // 9999
```

### Icon Size

```typescript
import { IconSize } from '@/theme';

IconSize.sm   // 16
IconSize.md   // 24
IconSize.lg   // 32
IconSize.xl   // 48
```

### Fonts

Platform-specific font families:

```typescript
import { Fonts } from '@/theme';

Fonts.sans     // System UI fonts
Fonts.serif    // Serif fonts
Fonts.rounded  // Rounded fonts
Fonts.mono     // Monospace fonts
```

## Common Patterns

### Screen Layout

```typescript
import { Surface, Text } from 'react-native-paper';
import { Spacing } from '@/theme';
import { StyleSheet } from 'react-native';

function MyScreen() {
  return (
    <Surface style={styles.container}>
      <Text variant="headlineLarge">Welcome</Text>
      <Text variant="bodyMedium">Content here</Text>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: Spacing.md,
  },
});
```

### Custom Card

```typescript
import { Card, Text, Button } from 'react-native-paper';
import { Spacing, BorderRadius } from '@/theme';
import { StyleSheet } from 'react-native';

function CustomCard() {
  return (
    <Card
      mode="elevated"
      style={[styles.card, { borderRadius: BorderRadius.lg }]}
    >
      <Card.Content style={styles.content}>
        <Text variant="titleLarge">Title</Text>
        <Text variant="bodyMedium">Description</Text>
      </Card.Content>
      <Card.Actions>
        <Button>Cancel</Button>
        <Button mode="contained">Confirm</Button>
      </Card.Actions>
    </Card>
  );
}

const styles = StyleSheet.create({
  card: {
    margin: Spacing.md,
  },
  content: {
    gap: Spacing.sm,
  },
});
```

### Conditional Styling with Theme

```typescript
import { useTheme } from 'react-native-paper';
import { Surface, Text } from 'react-native-paper';

function ThemedComponent() {
  const theme = useTheme();
  const isDark = theme.dark;

  return (
    <Surface style={{
      backgroundColor: isDark ? '#1a1a1a' : '#f5f5f5'
    }}>
      <Text style={{ color: theme.colors.onSurface }}>
        Adapts to theme!
      </Text>
    </Surface>
  );
}
```

## Customization

### Change App Colors

Edit `src/theme/colors.ts` to change your app's color scheme:

```typescript
export const Colors = {
  light: {
    tint: '#6200ee',      // Your brand color
    background: '#ffffff',
    text: '#000000',
    // ... other colors
  },
  dark: {
    tint: '#bb86fc',      // Dark mode brand color
    background: '#121212',
    text: '#ffffff',
    // ... other colors
  },
};
```

The Paper theme automatically picks up these changes!

### Extend Spacing/Metrics

Add new values in `src/theme/metrics.ts`:

```typescript
export const Spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
  huge: 64,     // Add custom size
};
```

### Add Custom Paper Theme Colors

Extend the Paper theme in `src/theme/paper-theme.ts`:

```typescript
export const paperLightTheme: MD3Theme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    primary: Colors.light.tint,
    // Add custom semantic colors
    success: '#10b981',
    warning: '#f59e0b',
    info: '#3b82f6',
  },
};
```

Then use TypeScript module augmentation if needed:

```typescript
// src/types/theme.d.ts
import { MD3Theme } from 'react-native-paper';

declare module 'react-native-paper' {
  interface MD3Colors {
    success: string;
    warning: string;
    info: string;
  }
}
```

## Best Practices

### ✅ Do's

1. **Always use Paper components** for UI (`Text`, `Button`, `Card`, etc.)
2. **Use `Spacing` constants** for all padding/margin values
3. **Access theme via `useTheme()`** for dynamic colors
4. **Define colors in one place** (`colors.ts`)
5. **Use Paper's text variants** (`titleLarge`, `bodyMedium`, etc.)

### ❌ Don'ts

1. **Don't hardcode spacing** values (use `Spacing` constants)
2. **Don't hardcode colors** (use Paper theme colors)
3. **Don't create custom themed wrappers** (Paper handles theming)
4. **Don't mix styling approaches** (stick to Paper + theme constants)

## Migration from Custom Components

If you see old code using `ThemedView` or `ThemedText`, migrate it:

```typescript
// ❌ Old (deprecated)
import { ThemedView } from '@/components/themed-view';
import { ThemedText } from '@/components/themed-text';

<ThemedView style={styles.container}>
  <ThemedText type="title">Hello</ThemedText>
</ThemedView>

// ✅ New (correct)
import { Surface, Text } from 'react-native-paper';

<Surface style={styles.container}>
  <Text variant="headlineMedium">Hello</Text>
</Surface>
```

## Resources

- [React Native Paper Documentation](https://callstack.github.io/react-native-paper/)
- [Material Design 3](https://m3.material.io/)
- [Paper Theme Reference](https://callstack.github.io/react-native-paper/docs/guides/theming)

