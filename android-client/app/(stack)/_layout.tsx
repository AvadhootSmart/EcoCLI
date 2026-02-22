import { Stack } from 'expo-router';
import { ThemedView } from '@/components/themed-view';

export default function StackLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: true,
        headerStyle: {
          backgroundColor: '#151718',
        },
        headerTintColor: '#fff',
        headerTitleStyle: {
          fontWeight: '600',
        },
        contentStyle: {
          backgroundColor: '#151718',
        },
      }}
    >
      <Stack.Screen
        name="index"
        options={{
          title: 'Eco',
        }}
      />
      <Stack.Screen
        name="logs"
        options={{
          title: 'Event Logs',
        }}
      />
      <Stack.Screen
        name="settings"
        options={{
          title: 'Settings',
        }}
      />
    </Stack>
  );
}
