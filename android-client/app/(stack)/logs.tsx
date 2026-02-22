import React from 'react';
import { StyleSheet, FlatList, View } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { useEco } from '@/context';
import type { LogEntry } from '@/types';

function LogItem({ log }: { log: LogEntry }) {
  const time = new Date(log.timestamp).toLocaleTimeString();
  const directionColor = log.direction === 'sent' ? '#4CAF50' : '#2196F3';

  return (
    <View style={styles.logItem}>
      <View style={styles.logHeader}>
        <ThemedText style={[styles.logDirection, { color: directionColor }]}>
          {log.direction.toUpperCase()}
        </ThemedText>
        <ThemedText style={styles.logType}>{log.type}</ThemedText>
        <ThemedText style={styles.logTime}>{time}</ThemedText>
      </View>
      <ThemedText style={styles.logPayload} numberOfLines={2}>
        {JSON.stringify(log.payload)}
      </ThemedText>
    </View>
  );
}

export default function LogsScreen() {
  const { logs } = useEco();

  return (
    <ThemedView style={styles.container}>
      <ThemedText type="title" style={styles.title}>Event Logs</ThemedText>
      
      {logs.length === 0 ? (
        <View style={styles.empty}>
          <ThemedText style={styles.emptyText}>No events yet</ThemedText>
        </View>
      ) : (
        <FlatList
          data={logs}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <LogItem log={item} />}
          contentContainerStyle={styles.list}
        />
      )}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
  },
  title: {
    marginBottom: 16,
  },
  list: {
    paddingBottom: 20,
  },
  logItem: {
    backgroundColor: 'rgba(255,255,255,0.05)',
    padding: 12,
    borderRadius: 8,
    marginBottom: 8,
  },
  logHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 4,
  },
  logDirection: {
    fontSize: 10,
    fontWeight: '700',
  },
  logType: {
    fontSize: 12,
    fontWeight: '600',
  },
  logTime: {
    fontSize: 10,
    opacity: 0.5,
    marginLeft: 'auto',
  },
  logPayload: {
    fontSize: 11,
    opacity: 0.7,
    fontFamily: 'monospace',
  },
  empty: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  emptyText: {
    opacity: 0.5,
  },
});
