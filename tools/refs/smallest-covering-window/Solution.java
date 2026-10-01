class Solution {
    public String smallestWindow(String text, String pattern) {
        int[] need = new int[128];
        for (char ch : pattern.toCharArray()) {
            need[ch]++;
        }
        int missing = pattern.length();
        int bestStart = 0;
        int bestLen = -1;
        int start = 0;
        for (int i = 0; i < text.length(); i++) {
            if (need[text.charAt(i)]-- > 0) {
                missing--;
            }
            while (missing == 0) {
                if (bestLen < 0 || i - start + 1 < bestLen) {
                    bestStart = start;
                    bestLen = i - start + 1;
                }
                if (++need[text.charAt(start++)] > 0) {
                    missing++;
                }
            }
        }
        return bestLen < 0 ? "" : text.substring(bestStart, bestStart + bestLen);
    }
}
