import { fetchHomeData } from '@/state/features/home/homeSlice';
import { useAppDispatch, useAppSelector } from '@/state/hooks';
import { Spacing } from '@/theme';
import { env } from '@/utils/env';
import { StyleSheet } from 'react-native';
import { Button, Card, Surface, Text } from 'react-native-paper';

export default function HomeScreen() {
  const dispatch = useAppDispatch();

  // Access feature-specific state from Redux
  const { data, isLoading, error } = useAppSelector((state) => state.home);

  // Access global state (example: auth)
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);

  const handleApiCall = () => {
    dispatch(fetchHomeData());
  };

  return (
    <Surface style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text variant="titleLarge">Home</Text>
          <Text variant="bodyMedium">
            {isAuthenticated && user
              ? `Welcome back, ${user.name}!`
              : 'Welcome to your app!'}
          </Text>

          <Text variant="bodyMedium" style={styles.meta}>
            Current environment: {env.appEnv}
          </Text>
          <Text variant="bodyMedium" style={styles.meta}>
            API base URL: {env.apiUrl}
          </Text>
          <Text variant="bodySmall" style={styles.meta}>
            💡 This screen uses Redux for state management
          </Text>

          <Button
            mode="contained"
            onPress={handleApiCall}
            loading={isLoading}
            disabled={isLoading}
            style={styles.button}
          >
            Call API
          </Button>

          {error && (
            <Card style={styles.resultCard} mode="outlined">
              <Card.Content>
                <Text variant="bodySmall" style={styles.errorText}>
                  Error: {error}
                </Text>
              </Card.Content>
            </Card>
          )}

          {data && !error && (
            <Card style={styles.resultCard} mode="outlined">
              <Card.Content>
                <Text variant="labelMedium" style={styles.successLabel}>
                  Success! ✓
                </Text>
                <Text variant="bodySmall">
                  {typeof data === 'string' ? data : JSON.stringify(data, null, 2)}
                </Text>
              </Card.Content>
            </Card>
          )}
        </Card.Content>
      </Card>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: Spacing.md,
    justifyContent: 'center',
    alignItems: 'center',
  },
  card: {
    width: '100%',
    maxWidth: 400,
  },
  meta: {
    marginTop: Spacing.sm,
  },
  button: {
    marginTop: Spacing.lg,
  },
  resultCard: {
    marginTop: Spacing.md,
  },
  errorText: {
    color: '#ef4444',
  },
  successLabel: {
    color: '#10b981',
    marginBottom: Spacing.sm,
  },
});
