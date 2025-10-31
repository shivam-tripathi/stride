import { Spacing } from '@/theme';
import { StyleSheet } from 'react-native';
import { Surface, Text } from 'react-native-paper';

export default function ModalScreen() {
  return (
    <Surface style={styles.container}>
      <Text variant="headlineMedium">Modal</Text>
      <Text variant="bodyMedium">This is a modal screen.</Text>
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: Spacing.lg,
  },
});
