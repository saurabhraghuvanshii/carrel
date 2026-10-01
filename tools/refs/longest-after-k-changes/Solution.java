class Solution {
    public int longestAfterChanges(String text, int k) {
        int[] count = new int[26];
        int best = 0;
        int most = 0;
        int start = 0;
        for (int i = 0; i < text.length(); i++) {
            most = Math.max(most, ++count[text.charAt(i) - 'A']);
            while (i - start + 1 - most > k) {
                count[text.charAt(start) - 'A']--;
                start++;
            }
            best = Math.max(best, i - start + 1);
        }
        return best;
    }
}
