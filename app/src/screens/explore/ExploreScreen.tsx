import { fetchExploreItems, setSearchQuery, setSelectedCategory } from '@/state/features/explore/exploreSlice';
import { useAppDispatch, useAppSelector } from '@/state/hooks';
import { Spacing } from '@/theme';
import { useEffect } from 'react';
import { FlatList, StyleSheet, View } from 'react-native';
import { ActivityIndicator, Card, Chip, Searchbar, Surface, Text } from 'react-native-paper';

export default function ExploreScreen() {
  const dispatch = useAppDispatch();

  // Access explore feature state from Redux
  const { items, isLoading, error, searchQuery, selectedCategory, page } = useAppSelector(
    (state) => state.explore
  );

  // Fetch initial data
  useEffect(() => {
    dispatch(fetchExploreItems({ page: 1, category: selectedCategory }));
  }, [dispatch, selectedCategory]);

  const handleSearch = (query: string) => {
    dispatch(setSearchQuery(query));
  };

  const handleCategorySelect = (category: string | null) => {
    dispatch(setSelectedCategory(category));
  };

  const renderEmptyState = () => {
    if (isLoading) {
      return (
        <View style={styles.centerContent}>
          <ActivityIndicator size="large" />
          <Text style={styles.loadingText}>Loading explore items...</Text>
        </View>
      );
    }

    if (error) {
      return (
        <View style={styles.centerContent}>
          <Text style={styles.errorText}>Error: {error}</Text>
          <Text style={styles.hintText}>
            This is a demo. The API endpoint may not be available yet.
          </Text>
        </View>
      );
    }

    return (
      <View style={styles.centerContent}>
        <Text style={styles.emptyText}>No items found</Text>
        <Text style={styles.hintText}>
          Try adjusting your filters or check back later
        </Text>
      </View>
    );
  };

  return (
    <Surface style={styles.container}>
      <View style={styles.header}>
        <Text variant="headlineMedium" style={styles.title}>Explore</Text>
        <Text variant="bodySmall" style={styles.subtitle}>
          💡 This screen uses Redux with pagination support
        </Text>
      </View>

      {/* Search Bar */}
      <Searchbar
        placeholder="Search explore items..."
        onChangeText={handleSearch}
        value={searchQuery}
        style={styles.searchBar}
      />

      {/* Category Filters */}
      <View style={styles.categoryContainer}>
        <Chip
          selected={selectedCategory === null}
          onPress={() => handleCategorySelect(null)}
          style={styles.chip}
        >
          All
        </Chip>
        <Chip
          selected={selectedCategory === 'tech'}
          onPress={() => handleCategorySelect('tech')}
          style={styles.chip}
        >
          Tech
        </Chip>
        <Chip
          selected={selectedCategory === 'design'}
          onPress={() => handleCategorySelect('design')}
          style={styles.chip}
        >
          Design
        </Chip>
        <Chip
          selected={selectedCategory === 'business'}
          onPress={() => handleCategorySelect('business')}
          style={styles.chip}
        >
          Business
        </Chip>
      </View>

      {/* Items List */}
      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <Card style={styles.card} mode="elevated">
            <Card.Content>
              <Text variant="titleMedium">{item.title}</Text>
              <Text variant="bodyMedium">{item.description}</Text>
            </Card.Content>
          </Card>
        )}
        ListEmptyComponent={renderEmptyState}
        contentContainerStyle={items.length === 0 ? styles.emptyList : undefined}
      />
    </Surface>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  header: {
    padding: Spacing.md,
    paddingBottom: Spacing.sm,
  },
  title: {
    fontWeight: 'bold',
  },
  subtitle: {
    marginTop: Spacing.xs,
    opacity: 0.7,
  },
  searchBar: {
    marginHorizontal: Spacing.md,
    marginBottom: Spacing.sm,
  },
  categoryContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: Spacing.md,
    paddingBottom: Spacing.sm,
  },
  chip: {
    marginRight: Spacing.sm,
    marginBottom: Spacing.sm,
  },
  card: {
    marginHorizontal: Spacing.md,
    marginBottom: Spacing.sm,
  },
  emptyList: {
    flex: 1,
  },
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: Spacing.xl,
  },
  loadingText: {
    marginTop: Spacing.md,
    opacity: 0.7,
  },
  errorText: {
    color: '#ef4444',
    marginBottom: Spacing.sm,
    textAlign: 'center',
  },
  emptyText: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: Spacing.sm,
    textAlign: 'center',
  },
  hintText: {
    opacity: 0.6,
    textAlign: 'center',
  },
});
