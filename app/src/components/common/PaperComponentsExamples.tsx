/**
 * Example components using React Native Paper
 * Use these as a reference for building your own components
 */

import { useState } from 'react';
import { StyleSheet, View } from 'react-native';
import {
    Button,
    Card,
    Chip,
    Dialog,
    FAB,
    Portal,
    Snackbar,
    Text,
    TextInput,
} from 'react-native-paper';

// Example 1: Simple Card with Actions
export function ExampleCard() {
  return (
    <Card mode="elevated" style={styles.card}>
      <Card.Cover source={{ uri: 'https://picsum.photos/700' }} />
      <Card.Title title="Card Title" subtitle="Card Subtitle" />
      <Card.Content>
        <Text variant="bodyMedium">
          This is an example card using React Native Paper components.
        </Text>
      </Card.Content>
      <Card.Actions>
        <Button>Cancel</Button>
        <Button mode="contained">OK</Button>
      </Card.Actions>
    </Card>
  );
}

// Example 2: Form with Text Inputs
export function ExampleForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  return (
    <View style={styles.form}>
      <TextInput
        label="Email"
        value={email}
        onChangeText={setEmail}
        mode="outlined"
        keyboardType="email-address"
        autoCapitalize="none"
        style={styles.input}
      />
      <TextInput
        label="Password"
        value={password}
        onChangeText={setPassword}
        mode="outlined"
        secureTextEntry
        style={styles.input}
      />
      <Button mode="contained" style={styles.button}>
        Sign In
      </Button>
    </View>
  );
}

// Example 3: Chips (Tags)
export function ExampleChips() {
  const [selected, setSelected] = useState<string[]>([]);

  const toggleChip = (value: string) => {
    setSelected((prev) =>
      prev.includes(value) ? prev.filter((v) => v !== value) : [...prev, value]
    );
  };

  return (
    <View style={styles.chipContainer}>
      <Chip
        selected={selected.includes('react')}
        onPress={() => toggleChip('react')}
        style={styles.chip}
      >
        React
      </Chip>
      <Chip
        selected={selected.includes('native')}
        onPress={() => toggleChip('native')}
        style={styles.chip}
      >
        Native
      </Chip>
      <Chip
        selected={selected.includes('paper')}
        onPress={() => toggleChip('paper')}
        style={styles.chip}
      >
        Paper
      </Chip>
    </View>
  );
}

// Example 4: Dialog
export function ExampleDialog() {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <Button onPress={() => setVisible(true)}>Show Dialog</Button>
      <Portal>
        <Dialog visible={visible} onDismiss={() => setVisible(false)}>
          <Dialog.Title>Alert</Dialog.Title>
          <Dialog.Content>
            <Text variant="bodyMedium">
              This is an example dialog using React Native Paper.
            </Text>
          </Dialog.Content>
          <Dialog.Actions>
            <Button onPress={() => setVisible(false)}>Cancel</Button>
            <Button onPress={() => setVisible(false)}>OK</Button>
          </Dialog.Actions>
        </Dialog>
      </Portal>
    </>
  );
}

// Example 5: Snackbar (Toast notification)
export function ExampleSnackbar() {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <Button onPress={() => setVisible(true)}>Show Snackbar</Button>
      <Snackbar
        visible={visible}
        onDismiss={() => setVisible(false)}
        duration={3000}
        action={{
          label: 'Undo',
          onPress: () => {
            // Handle undo
          },
        }}
      >
        Item deleted successfully
      </Snackbar>
    </>
  );
}

// Example 6: Floating Action Button
export function ExampleFAB() {
  return (
    <FAB
      icon="plus"
      style={styles.fab}
      onPress={() => console.log('FAB pressed')}
    />
  );
}

const styles = StyleSheet.create({
  card: {
    margin: 16,
  },
  form: {
    padding: 16,
  },
  input: {
    marginBottom: 12,
  },
  button: {
    marginTop: 8,
  },
  chipContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: 16,
  },
  chip: {
    margin: 4,
  },
  fab: {
    position: 'absolute',
    margin: 16,
    right: 0,
    bottom: 0,
  },
});

