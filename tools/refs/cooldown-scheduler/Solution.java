class Solution {
    public int fewestUnits(String tasks, int n) {
        int[] count = new int[26];
        int most = 0;
        for (char ch : tasks.toCharArray()) {
            most = Math.max(most, ++count[ch - 'A']);
        }
        int withMost = 0;
        for (int c : count) {
            if (c == most) {
                withMost++;
            }
        }
        return Math.max(tasks.length(), (most - 1) * (n + 1) + withMost);
    }
}
