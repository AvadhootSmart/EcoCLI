import React from 'react';
import { StyleSheet, View, Switch, TouchableOpacity } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { useEco } from '@/context';

export default function SettingsScreen() {
  const { state, deviceId, serverUrl } = useEco();

  return (
    <ThemedView style={styles.container}>
      <ThemedText type="title" style={styles.title}>Settings</ThemedText>

      <ThemedView style={styles.section}>
        <ThemedText type="subtitle">Device Info</ThemedText>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Device ID</ThemedText>
          <ThemedText style={styles.value}>{deviceId || 'Not set'}</ThemedText>
        </View>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Server</ThemedText>
          <ThemedText style={styles.value}>{serverUrl}</ThemedText>
        </View>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Connection</ThemedText>
          <ThemedText style={styles.value}>{state}</ThemedText>
        </View>
      </ThemedView>

      <ThemedView style={styles.section}>
        <ThemedText type="subtitle">Permissions</ThemedText>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Clipboard</ThemedText>
          <Switch value={true} onValueChange={() => {}} />
        </View>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Notifications</ThemedText>
          <Switch value={true} onValueChange={() => {}} />
        </View>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Phone State</ThemedText>
          <Switch value={true} onValueChange={() => {}} />
        </View>
      </ThemedView>

      <ThemedView style={styles.section}>
        <ThemedText type="subtitle">About</ThemedText>
        <View style={styles.row}>
          <ThemedText style={styles.label}>Version</ThemedText>
          <ThemedText style={styles.value}>1.0.0</ThemedText>
        </View>
      </ThemedView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
  },
  title: {
    marginBottom: 24,
  },
  section: {
    marginBottom: 24,
    padding: 16,
    borderRadius: 12,
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: 'rgba(255,255,255,0.1)',
  },
  label: {
    fontSize: 14,
  },
  value: {
    fontSize: 14,
    opacity: 0.7,
  },
});
