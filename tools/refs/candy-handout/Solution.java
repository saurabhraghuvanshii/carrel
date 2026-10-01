class Solution {
    public int fewestSweets(int[] ratings) {
        int n = ratings.length;
        int[] sweets = new int[n];
        for (int i = 0; i < n; i++) {
            sweets[i] = i > 0 && ratings[i] > ratings[i - 1] ? sweets[i - 1] + 1 : 1;
        }
        int total = sweets[n - 1];
        for (int i = n - 2; i >= 0; i--) {
            if (ratings[i] > ratings[i + 1]) {
                sweets[i] = Math.max(sweets[i], sweets[i + 1] + 1);
            }
            total += sweets[i];
        }
        return total;
    }
}
