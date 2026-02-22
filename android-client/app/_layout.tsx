import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import 'react-native-reanimated';

import { useColorScheme } from '@/hooks/use-color-scheme';
import { EcoProvider } from '@/context';

export const unstable_settings = {
  anchor: '(stack)',
};

export default function RootLayout() {
  const colorScheme = useColorScheme();

  return (
    <ThemeProvider value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}>
      <EcoProvider>
        <Stack>
          <Stack.Screen name="(stack)" options={{ headerShown: false }} />
        </Stack>
      </EcoProvider>
      <StatusBar style="auto" />
    </ThemeProvider>
  );
}
